package module

import (
	"log"
	"net/url"
	"time"
)

type Message struct {
	Protocol string    `json:"protocol"`
	Data     string    `json:"data"`
	Module   string    `json:"module"`
	Perms    int       `json:"perms"`
	Time     time.Time `json:"time"`
}

type Module struct {
	Server         string
	Name, Password string
	Protocols      []string

	Refresh time.Duration

	HSet Handler

	quit chan bool
}

func NewModule(server, name, password string) (*Module, error) {

}

func (m *Module) Close() {
	m.quit <- true
}

func (m *Module) Listen() {
	m.quit = make(chan bool)

	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(time.Second):
				m.Poll()
			}
		}
	}()
}

func (m *Module) Poll() {
	u, err := url.Parse(m.Server)
	if err != nil {
		log.Panic("Invalid Server URL:", m.Server)
	}

	u.RawPath = "/get"

}
