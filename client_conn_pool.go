package gosocket

import (
	"net"
	"unsafe"
	"sync/atomic"
)

// ClientConnPool manages a pool of client connections.
type ClientConnPool interface {
	// GetConn gets a idle connection from the pool.
	GetConn() (Conn, error)
	// Close closes all idle connections.
	Close() error
}

type clientConnPool struct {
	addr         string
	connCreator  func(nc net.Conn) Conn
	connAddr     *unsafe.Pointer
	dialer       MyDialer
	dialCallAddr *unsafe.Pointer
}

// MyDialer defines how to connect to the address on the named network.
type MyDialer interface {
	Dial(network, address string) (net.Conn, error)
}

type dialCall struct {
	// done is closed when dialing process is completed.
	done chan struct{}
	// resp is valid after channel 'done' is closed.
	resp Conn
	// err is valid after channel 'done' is closed
	err error
}

func newClientConnPool(addr string, dialer MyDialer, connCreator func(nc net.Conn) Conn) *clientConnPool {
	p := &clientConnPool{
		addr:        addr,
		connCreator: connCreator,
		dialer:      dialer,
	}
	var addr1 *byte
	p.connAddr = (*unsafe.Pointer)(unsafe.Pointer(&addr1))
	var addr2 *byte
	p.dialCallAddr = (*unsafe.Pointer)(unsafe.Pointer(&addr2))
	return p
}

func (p *clientConnPool) GetConn() (Conn, error) {
	conn, ok := p.getConn()
	if !ok {
		return p.doDial()
	}
	return conn, nil
}

func (p *clientConnPool) Close() error{
	conn,ok:=p.getConn()
	if !ok{
		return nil
	}
	return conn.Close()
}


func (p *clientConnPool) doDial() (Conn, error) {
	call := &dialCall{done: make(chan struct{})}
	dialing := (*dialCall)(atomic.LoadPointer(p.dialCallAddr))
	swapped := atomic.CompareAndSwapPointer(p.dialCallAddr, nil, unsafe.Pointer(call))
	if !swapped {
		// A dial call is already in-flight. Don't start another.
		<-dialing.done
		return dialing.resp, dialing.err
	}
	// double check.
	if conn, ok := p.getConn(); ok {
		return conn,nil
	}
	// Start to connect to the remote server.
	var nc net.Conn
	nc, call.err = p.dialer.Dial("tcp", p.addr)
	if call.err==nil{
		call.resp = p.connCreator(nc)
		p.putConn(call.resp)
	}
	// Note: Deleting the dial call must be performed after storing the connection
	// If the address is successfully dialed.
	atomic.StorePointer(p.dialCallAddr, nil)
	close(call.done)
	return call.resp, call.err
}

// getInvalidConn returns the connection that can take a new request, (i.e., (respCh, true)).
// If the connection does not exist it returns (nil,false).
func (p *clientConnPool) getConn() (Conn, bool) {
	val := atomic.LoadPointer(p.connAddr)
	if val == nil {
		return nil, false
	}
	conn := *(*Conn)(val)
	return conn, true
}

// putConn stores the connection.
func (p *clientConnPool) putConn(conn Conn) {
	atomic.StorePointer(p.connAddr, unsafe.Pointer(&conn))
}
