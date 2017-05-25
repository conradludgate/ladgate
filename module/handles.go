package module

import "sync"

type Handler interface {
	ServeMessage(Message)
}

type HandlerFunc func(Message)

func NotFound(message Message) {}
func NotFoundHandler = HandlerFunc(NotFound)

func (f HandlerFunc) ServeMessage(message Message) {
	f(message)
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

func (hs *HandlerSet) Handler(message Message) (h Handler) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	h, ok := hs.h[message.Protocol]

	if !ok {
		return NotFoundHandler
	}

	return h
}

func (hs *HandlerSet) ServeMessage(message Message) {
	h, _ := hs.Handler(message)
	h.Handle(message)
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

func (hs *HandlerSet) HandleFunc(protocol string, handler func(ResponseWriter, *Request)) {
	hs.Handle(pattern, HandlerFunc(handler))
}

func Handle(protocol string, handler Handler) { DefaultHandlerSet.Handle(protocol, handler) }
func HandleFunc(protocol string, handler func(ResponseWriter, *Request)) {
	DefaultHandlerSet.Handle(pattern, HandlerFunc(handler))
}