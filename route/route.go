package route

import (
	"sync"
	"reflect"
	"github.com/happyxcj/gosocket"
	"github.com/happyxcj/gosocket/protocol"
	"github.com/happyxcj/gosocket/pkts"
	"fmt"
)

var globalRC *RouterCenter

var ctxPool sync.Pool

func init() {
	globalRC = NewRouterCenter()
	ctxPool.New = func() interface{} {
		return new(Context)
	}
}

func Group(kind protocol.PktKind, handlers ...HandlerFunc) *RouterGroup {
	return globalRC.Group(kind, handlers...)
}

func Use(handlers ...HandlerFunc) *RouterCenter {
	return globalRC.Use(handlers...)
}

func UseAfter(handlers ...HandlerFunc) *RouterCenter {
	return globalRC.UseAfter(handlers...)
}

func Get(kind protocol.PktKind, msgId string) (*RouterInfo, bool) {
	return globalRC.Get(kind, msgId)
}

func GetContext() *Context {
	ctx := ctxPool.Get().(*Context)
	return ctx
}

func PutContext(ctx *Context) {
	ctxPool.Put(ctx)
}

// HandlePacket handle the given p for the c.
// It Returns a bool indicates whether the p can be handled successfully.
func HandlePacket(c gosocket.Conn, p protocol.Packet) bool {
	info, ok := Get(p.Kind(), GenMsgId(p))
	if !ok {
		return false
	}
	var msg interface{}
	if info.msgType != nil {
		msg = reflect.New(info.msgType).Interface()
	}
	ctx := GetContext()
	ctx.Reset(c, p, msg, info.handlers)
	ctx.Next()
	PutContext(ctx)
	return true
}

// GenMsgId returns a unique id of the packet application message.
func GenMsgId(p protocol.Packet) string {
	if p.Kind() == pkts.KindPing {
		return "ping"
	}
	pkt := p.(pkts.DataPkt)
	if pkt.Version() == 0 && pkt.Codec() == 0 {
		return fmt.Sprint(pkt.Cmd())
	}
	return fmt.Sprintf("%v-%v-%v", pkt.Cmd(), pkt.Version(), pkt.Codec())
}