package gosocket

import (
	"sync/atomic"
	"github.com/happyxcj/gosocket/protocol"
	"fmt"
	"errors"
)

const (
	defaultSendChSize = 200
)

var (
	// ErrClosedActively signals the connection is closed actively by the local peer.
	ErrClosedActively = errors.New("local peer actively closes the session")

	// ErrConnClosed is returned when we attempt to send a packet to a closed connection.
	ErrConnClosed = errors.New("the connection has been closed")

	// ErrSlowConsumer is returned when we attempt to enqueue a packet to a full sending channel.
	ErrSlowConsumer = errors.New("slow consumer detected")
)

// autoId is used to generate connection unique id.
var autoId uint64

var _ Conn = (*QueueConn)(nil)

type QueueConn struct {
	Conn

	// id is the unique connection id.
	id uint64

	// sendCh is the channel for send packet.
	// A zero value of it's size means disable a send channel.
	sendCh chan protocol.Packet

	sendChSize int

	// doneFlag indicates whether both the sending goroutine and receiving goroutine have quit.
	doneFlag uint32

	// closedFlag indicates whether the connection is closed.
	closedFlag uint32

	// cause represents the reason why the connection was closed.
	cause error

	// pktHandler handles the every received packet.
	pktHandler func(p protocol.Packet)

	// pendingPktsHandler handles the packets to be sent when both the sending goroutine
	// and receiving goroutine have quit.
	pendingPktsHandler func(pkts []protocol.Packet)

	// onClose is the callback when both the sending goroutine and receiving goroutine have quit.
	// The cause represents the reason why the connection was closed.
	onClose func(cause error)
}

type QueueConnOpt func(*QueueConn)

// SendChSize returns a QueueConnOpt to set sending channel size for the QueueConn.
func SendChSize(size int) QueueConnOpt {
	return func(c *QueueConn) {
		c.sendChSize = size
	}
}

// PktHandler returns a QueueConnOpt to set the packet handler for the QueueConn.
func PktHandler(handler func(protocol.Packet)) QueueConnOpt {
	return func(c *QueueConn) {
		c.pktHandler = handler
	}
}

// PendingPktsHandler returns a QueueConnOpt to set the pending packets handler for the QueueConn.
func PendingPktsHandler(handler func([]protocol.Packet)) QueueConnOpt {
	return func(c *QueueConn) {
		c.pendingPktsHandler = handler
	}
}

// OnClose returns a QueueConnOpt to set the closing callback for the QueueConn.
func OnClose(onClose func(error)) QueueConnOpt {
	return func(c *QueueConn) {
		c.onClose = onClose
	}
}

func NewQueueConn(conn Conn, opts ...QueueConnOpt) *QueueConn {
	c := &QueueConn{
		Conn:       conn,
		id:         atomic.AddUint64(&autoId, 1),
		sendChSize: defaultSendChSize,
		pktHandler: func(p protocol.Packet) {
			fmt.Println("receive packet: ", p.Desc())
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.sendCh = make(chan protocol.Packet, c.sendChSize)
	go c.sendLoop()
	go c.receiveLoop()
	return c
}

// PktHandler returns the packet handler.
func (c *QueueConn) PktHandler() func(p protocol.Packet) {
	return c.pktHandler
}

// Id returns the unique id of the connection.
func (c *QueueConn) Id() uint64 {
	return c.id
}

// IsClosed returns a bool indicating whether the connection has been closed.
func (c *QueueConn) IsClosed() bool {
	return atomic.LoadUint32(&c.closedFlag) != 0
}

// Enqueue attempts to send a packet to the sending channel, it will returns an error
// if the connection has been closed or the channel buffer is full.
// Be careful that it will close the connection right away if the channel buffer is full.
func (c *QueueConn) Send(p protocol.Packet) error {
	if c.IsClosed() {
		// Can't send message to the closed connection.
		return ErrConnClosed
	}
	select {
	case c.sendCh <- p:
		return nil
	default:
		// A slow consumer was detected, close the underlying connection right away.
		c.Conn.Close()
		// Mark the connection as closed.
		c.markClosed()
		return ErrSlowConsumer
	}
}

// Close closes the connection after all pending packets in sender queue are sent.
func (c *QueueConn) Close() error {
	err := c.Send(nil)
	// Mark the connection as closed.
	c.markClosed()
	return err
}

// CloseWith closes the connection after the given p is sent.
func (c *QueueConn) CloseWith(p protocol.Packet) error {
	c.Send(p)
	return c.Close()
}

func (c *QueueConn) sendLoop() {
	var err error
	for {
		p := <-c.sendCh
		if p == nil {
			break
		}
		err = c.Conn.Send(p)
		if err != nil {
			break
		}
	}
	if err == nil {
		err = ErrClosedActively
	}
	c.close(err, false)
}

func (c *QueueConn) receiveLoop() () {
	var err error
	for {
		var p protocol.Packet
		p, err = c.Conn.Receive()
		if err != nil {
			break
		}
		// handle message
		c.pktHandler(p)
	}
	c.close(err, true)
}

// markClosed only marks the connection status as closed,
// it will not close the underlying connection.
func (c *QueueConn) markClosed() {
	atomic.CompareAndSwapUint32(&c.closedFlag, 0, 1)
}

// closeWithErr closes the connection and handle the closing callback
// by the given err.
func (c *QueueConn) close(err error, isReadingErr bool) {
	if atomic.CompareAndSwapUint32(&c.doneFlag, 0, 1) {
		// Just the first error is the real cause.
		c.cause = err
		// Mark the connection as closed.
		c.markClosed()
		// Close the underlying connection right away.
		// Situation 1:
		// 		It will trigger the receiving goroutine to exit when the sending goroutine quits.
		// Situation 2:
		// 		It will trigger the sending goroutine to exit if the sending channel is not empty
		// 		when the receiving goroutine quits.
		c.Conn.Close()

		if isReadingErr {
			// Ensure to notify the sending goroutine to exit when
			// the sending channel is empty.
			select {
			case c.sendCh <- nil:
			default:
			}
		}
		return
	}
	// Both the sending goroutine and reading goroutine have exit.
	if c.pendingPktsHandler != nil && len(c.sendCh) > 0 {
		c.pendingPktsHandler(c.getPendingPackets())
	}
	if c.onClose != nil {
		c.onClose(c.cause)
	}
}

func (c *QueueConn) getPendingPackets() []protocol.Packet {
	pkts := make([]protocol.Packet, 0, len(c.sendCh))
	for {
		select {
		case pkt := <-c.sendCh:
			if pkt == nil {
				return pkts
			}
			pkts = append(pkts, pkt)
		default:
			return pkts
		}
	}
}
