package module

import (
	"sync"
)

type Handler interface {
	ServeMessage(*Module, Message)
}

type HandlerFunc func(*Module, Message)

func NotFound(m *Module, message Message) {}

var NotFoundHandler = HandlerFunc(NotFound)

func (f HandlerFunc) ServeMessage(m *Module, message Message) {
	f(m, message)
}

type HandlerSet struct {
	mu sync.RWMutex
	h  map[string]handlerEntry
}

func NewHandlerSet() *HandlerSet { return new(HandlerSet) }

var DefaultHandlerSet = NewHandlerSet()

type handlerEntry struct {
	explicit bool
	h        Handler
	protocol string
}

func (hs *HandlerSet) Handler(message Message) Handler {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	h, ok := hs.h[message.Protocol]

	if !ok {
		return NotFoundHandler
	}

	return h.h
}

func (hs *HandlerSet) ServeMessage(m *Module, message Message) {
	hs.Handler(message).ServeMessage(m, message)
}

func (hs *HandlerSet) Handle(protocol string, handler Handler) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if protocol == "" {
		panic("module: invalid protocol " + protocol)
	}
	if handler == nil {
		panic("module: nil handler")
	}

	if hs.h[protocol].explicit {
		panic("module: multiple registrations for " + protocol)
	}

	if hs.h == nil {
		hs.h = make(map[string]handlerEntry)
	}
	hs.h[protocol] = handlerEntry{explicit: true, h: handler, protocol: protocol}
}

func (hs *HandlerSet) HandleFunc(protocol string, handler func(*Module, Message)) {
	hs.Handle(protocol, HandlerFunc(handler))
}

func Handle(protocol string, handler Handler) { DefaultHandlerSet.Handle(protocol, handler) }
func HandleFunc(protocol string, handler func(*Module, Message)) {
	DefaultHandlerSet.Handle(protocol, HandlerFunc(handler))
}
