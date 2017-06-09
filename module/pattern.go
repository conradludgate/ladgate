package module

import "sync"

type Pattern interface {
	Match(string) bool
}

type PatternFunc func(string) bool

func (pf PatternFunc) Match(in string) bool {
	return pf(in)
}

type PatternHandler struct {
	mu sync.RWMutex
	h  []*PatternHandle
}

type PatternHandle struct {
	h       Handler
	pattern Pattern

	perms     int
	exclusive bool
}

func NewPatternHandler() *PatternHandler { return new(PatternHandler) }

func (ph *PatternHandle) Perms(n int) *PatternHandle {
	if n == 0 {
		panic("module: no module will satisfy perms of 0")
	}

	ph.perms = n
	return ph
}

func (ph *PatternHandle) Exclusive() *PatternHandle {
	ph.exclusive = true
	return ph
}

var DefaultPatternHandler = NewPatternHandler()

func (ph *PatternHandler) ServeMessage(m *Module, msg Message) {
	for _, v := range ph.h {
		if v.perms&msg.Perms > 0 && v.pattern.Match(msg.Data) {
			v.h.ServeMessage(m, msg)

			if v.exclusive {
				return
			}
		}
	}
}

func (ph *PatternHandler) Handle(pattern Pattern, handler Handler) *PatternHandle {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	if pattern == nil {
		panic("module: nil pattern")
	}

	if handler == nil {
		panic("module: nil handler")
	}

	handle := PatternHandle{handler, pattern, -1, false}
	ph.h = append(ph.h, &handle)

	return &handle
}

func (ph *PatternHandler) HandleFunc(pattern Pattern, handler func(*Module, Message)) {
	ph.Handle(pattern, HandlerFunc(handler))
}
