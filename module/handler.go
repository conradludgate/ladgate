package module

type Handler interface {
	ServeMessage(*Module, Message)
}

type HandlerFunc func(*Module, Message)

func (f HandlerFunc) ServeMessage(m *Module, message Message) {
	f(m, message)
}

func NotFound(m *Module, message Message) {}

var NotFoundHandler = HandlerFunc(NotFound)

func Handle(pattern Pattern, handler Handler) *PatternHandle {
	return DefaultPatternHandler.Handle(pattern, handler)
}

func HandleFunc(pattern Pattern, handler func(*Module, Message)) *PatternHandle {
	return DefaultPatternHandler.Handle(pattern, HandlerFunc(handler))
}
