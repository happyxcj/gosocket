package gosocket

import (
	"github.com/happyxcj/gosocket/protocol"
	"net"
	"time"
)

const (
	defaultWriteTimeout = 15 * time.Second
	defaultReadTimeout  = 10 * time.Second
)


var _ Conn = (*EasyConn)(nil)

type EasyConn struct {
	nc    net.Conn
	codec *protocol.Codec
	// writeTimeout represents the deadline duration for future write calls
	// and any currently-blocked Read call.
	writeTimeout time.Duration
	// readTimeout represents the deadline duration for future Read calls
	// and any currently-blocked Read call.
	readTimeout time.Duration
}

type EasyConnOpt func(*EasyConn)

// WriteTimeout returns a EasyConnOpt to set the writing timeout for the connection.
func WriteTimeout(timeout time.Duration) EasyConnOpt {
	return func(c *EasyConn) {
		c.writeTimeout = timeout
	}
}

// ReadTimeout returns a EasyConnOpt to set the reading timeout for the connection.
func ReadTimeout(timeout time.Duration) EasyConnOpt {
	return func(c *EasyConn) {
		c.writeTimeout = timeout
	}
}

// ReadTimeout returns a EasyConnOpt to set the packet codec for the connection.
func Codec(codec *protocol.Codec) EasyConnOpt {
	return func(c *EasyConn) {
		c.codec = codec
	}
}

func NewEasyConn(nc net.Conn, opts ...EasyConnOpt) *EasyConn {
	c := &EasyConn{
		nc: nc,
		codec: protocol.NewCodec(protocol.NewWriter(nc, protocol.BigEndian),
			protocol.NewReader(nc, protocol.BigEndian)),
		writeTimeout: defaultWriteTimeout,
		readTimeout:  defaultReadTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// UnderlyingConn returns the internal net.Conn.
func (c *EasyConn) UnderlyingConn() net.Conn {
	return c.nc
}

func (c *EasyConn) Send(p protocol.Packet) error {
	if c.writeTimeout > 0 {
		c.nc.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	return c.codec.Write(p)
}

func (c *EasyConn) Receive() (protocol.Packet, error) {
	if c.readTimeout > 0 {
		c.nc.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	return c.codec.Read()
}

func (c *EasyConn) Close() error {
	return c.nc.Close()
}

