package route

import (
	"math"
	"fmt"
	"github.com/happyxcj/gosocket"
	"github.com/happyxcj/gosocket/protocol"
)

const abortIndex int8 = math.MaxInt8 / 2

type Context struct {
	Conn     gosocket.Conn
	Pkt      protocol.Packet
	Msg      interface{}
	index    int8
	handlers []HandlerFunc
	// values is a key/value pair exclusively for the context of each request.
	//
	// It it usually used to store a value first using c.Set method so that this context
	// can read the value later using c.Get or c.MustGet method.
	values map[Key]interface{}
}

type HandlerFunc func(*Context)

type Key int

func NewContext() *Context {
	ctx := &Context{index: -1}
	return ctx
}

func (c *Context) Reset(conn gosocket.Conn, pkt protocol.Packet, msg interface{}, handlers []HandlerFunc) {
	c.Conn = conn
	c.Pkt = pkt
	c.Msg = msg
	c.handlers = handlers
	c.index = -1
	c.values = nil
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents the pending handlers from being called.
// Note that this will not stop the current handler.
func (c *Context) Abort() {
	c.index = abortIndex
}

// SkipNext prevents the next pending handler from being called.
// Note that this will not stop the current handler.
func (c *Context) SkipNext() {
	c.index += 1
}

// SkipNext prevents the next 'n' pending handlers from being called.
// Note that this will not stop the current handler.
func (c *Context) SkipNextN(n int) {
	c.index += int8(n)
}

// Next executes the pending handlers in the handler chain.
func (c *Context) Next() {
	c.index++
	for n := int8(len(c.handlers)); c.index < n; c.index++ {
		c.handlers[c.index](c)
	}
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes c.values if it was not used previously.
func (c *Context) Set(key Key, v interface{}) {
	if c.values == nil {
		c.values = make(map[Key]interface{})
	}
	c.values[key] = v
}

// Get returns the value for the given key (i.e., (value, true)).
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key Key) (interface{}, bool) {
	v, exist := c.values[key]
	return v, exist
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key Key) interface{} {
	if v, ok := c.Get(key); ok {
		return v
	}
	panic(fmt.Sprintf("key %v does not exist", key))
}
