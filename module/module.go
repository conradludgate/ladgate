package module

import "time"

type Message struct {
	Protocol string    `json:"protocol"`
	Data     string    `json:"data"`
	Module   string    `json:"module"`
	Perms    int       `json:"perms"`
	Time     time.Time `json:"time"`
}

type Module struct {
	Server string

	Refresh time.Duration

	Name, Password string

	HSet Handler
}

func NewModule(server, name, password string, handler Handler) *Module {
	if handler == nil {
		handler = DefaultHandlerSet
	}

	return &Module{
		server,
		100 * time.Millisecond,
		name,
		password,
		handler,
	}
}

func (m *Module) Poll() {

}
