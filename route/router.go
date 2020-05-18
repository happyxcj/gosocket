package route

import (
	"reflect"
	"fmt"
	"github.com/happyxcj/gosocket/protocol"
)

type RouterCenter struct {
	handlers []HandlerFunc
	headsNum int
	routers  map[protocol.PktKind]*RouterGroup
}

func NewRouterCenter() *RouterCenter {
	return &RouterCenter{
		routers: make(map[protocol.PktKind]*RouterGroup),
	}
}

func (rc *RouterCenter) Use(handlers ...HandlerFunc) *RouterCenter {
	heads := rc.handlers[:rc.headsNum]
	tails := rc.handlers[rc.headsNum:]
	rc.handlers = append(heads, handlers...)
	rc.handlers = append(rc.handlers, tails...)
	// Reset the number of head handlers.
	rc.headsNum = len(heads) + len(handlers)
	return rc
}

func (rc *RouterCenter) UseAfter(handlers ...HandlerFunc) *RouterCenter {
	rc.handlers = append(rc.handlers, handlers...)
	return rc
}

func (rc *RouterCenter) Group(kind protocol.PktKind, handlers ...HandlerFunc) *RouterGroup {
	handlers = combineHandlers(rc.handlers, rc.headsNum, handlers...)
	rg := &RouterGroup{
		handlers: handlers,
		headsNum: len(handlers),
		kind:     kind,
	}
	rc.routers[kind] = rg
	return rg
}

func (rc *RouterCenter) Get(kind protocol.PktKind, msgId string) (*RouterInfo, bool) {
	router, ok := rc.routers[kind]
	if !ok {
		return nil, false
	}
	info, ok := router.routes[msgId]
	return info, ok
}

type RouterGroup struct {
	kind     protocol.PktKind
	handlers []HandlerFunc
	headsNum int
	routes   map[string]*RouterInfo
}

type RouterInfo struct {
	msgType  reflect.Type
	handlers []HandlerFunc
}

func (rg *RouterGroup) Group() *RouterGroup {
	rg.handlers = nil
	rg.headsNum = 0
	return rg
}

func (rg *RouterGroup) Use(handlers ...HandlerFunc) *RouterGroup {
	heads := rg.handlers[:rg.headsNum]
	tails := rg.handlers[rg.headsNum:]
	rg.handlers = append(heads, handlers...)
	rg.handlers = append(rg.handlers, tails...)
	// Reset the number of head handlers.
	rg.headsNum = len(heads) + len(handlers)
	return rg
}

func (rg *RouterGroup) UseAfter(handlers ...HandlerFunc) *RouterGroup {
	rg.handlers = append(rg.handlers, handlers...)
	return rg
}

func (rg *RouterGroup) Handle(msgId, msg interface{}, handlers ...HandlerFunc) *RouterGroup {
	handlers = combineHandlers(rg.handlers, rg.headsNum, handlers...)
	if rg.routes == nil {
		rg.routes = make(map[string]*RouterInfo)
	}
	rg.routes[fmt.Sprint(msgId)] = &RouterInfo{
		handlers: handlers,
		msgType:  reflect.TypeOf(msg),
	}
	return rg
}

func combineHandlers(src []HandlerFunc, headsNum int, handlers ...HandlerFunc) []HandlerFunc {
	finalSize := len(src) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, src[:headsNum])
	copy(mergedHandlers[headsNum:], handlers)
	copy(mergedHandlers[headsNum+len(handlers):], src[headsNum:])
	return mergedHandlers
}
