
package gosocket

import (
	"time"
	"net"
	"github.com/happyxcj/gosocket/protocol"
)

const (
	defaultDialTimeout = 5 * time.Second
)

// Client represents a client connection to an specified server.
type Client struct {
	// opts contains the options to dial a server.
	opts      *DialOpts
	pool      *clientConnPool
}

type DialOpts struct {
	ecOpts      []EasyConnOpt
	qcOpts      []QueueConnOpt
	dialTimeout time.Duration
	// respTimeout specifies the timeout to wait for a server's response.
	// It's default value is "10*time.second".
	//respTimeout time.Duration
	dialer      MyDialer
}

// DialOpt specifies an option for a connection.
type DialOpt func(*DialOpts)

// ECOption returns a DialOpt to add a option to the internal ecOpts.
func ECOptions(opts ...EasyConnOpt) DialOpt {
	return func(o *DialOpts) {
		o.ecOpts = append(o.ecOpts, opts...)
	}
}

// QCOption returns a DialOpt to add a option to the internal qcOpts.
func QCOptions(opts ...QueueConnOpt) DialOpt {
	return func(o *DialOpts) {
		o.qcOpts = append(o.qcOpts, opts...)
	}
}

// RespTimeout returns a DialOpt to set the timeout to wait for a server's response.
func RespTimeout(t time.Duration) DialOpt {
	return func(o *DialOpts) {
		o.dialTimeout = t
	}
}

// DialTimeout returns a DialOpt to set the dial timeout for connecting to the server.
func DialTimeout(t time.Duration) DialOpt {
	return func(o *DialOpts) {
		o.dialTimeout = t
	}
}

// Dialer returns a DialOpt to set the dialer used to connect to the server.
func Dialer(dialer MyDialer) DialOpt {
	return func(o *DialOpts) {
		o.dialer = dialer
	}
}

func NewClient(addr string, opts ... DialOpt) *Client {
	c := &Client{
	}
	c.opts = &DialOpts{
		dialTimeout: defaultDialTimeout,
		dialer:      &net.Dialer{Timeout: defaultDialTimeout},
	}
	for _, opt := range opts {
		opt(c.opts)
	}
	connCreator := func(nc net.Conn) Conn {
		inner := NewEasyConn(nc, c.opts.ecOpts...)
		conn := NewQueueConn(inner, c.opts.qcOpts...)
		return conn
	}
	c.pool = newClientConnPool(addr, c.opts.dialer, connCreator)
	return c
}

// NewAndInitClient returns a Client and initializes to connect to the server.
func NewAndInitClient(addr string, opts ... DialOpt) (*Client, error) {
	c := NewClient(addr, opts...)
	_, err := c.pool.GetConn()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Send sends a packet to the server.
func (c *Client) Send(p protocol.Packet) error {
	conn, err := c.pool.GetConn()
	if err != nil {
		return err
	}
	return conn.Send(p)
}

// Close closes all idle connections in the pool.
func (c *Client) Close() error  {
	return c.pool.Close()
}
