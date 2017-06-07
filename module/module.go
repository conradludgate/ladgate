package module

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
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

	HSet Handler

	conns map[string]*websocket.Conn
}

func NewModule(server, name, password string) *Module {
	var m Module
	m.Server = server
	m.Name = name
	m.Password = password

	m.conns = make(map[string]*websocket.Conn)

	return &m
}

type wsInit struct {
	Protocol string `json:"protocol"`
	Module   string `json:"module"`
	Password string `json:"password"`
}

func (m *Module) Listen(protocol string, h Handler) error {
	if h == nil {
		m.HSet = DefaultHandlerSet
	} else {
		m.HSet = h
	}

	u, err := url.Parse(m.Server)

	if err != nil {
		return err
	}

	u.Path = "ws"

	m.conns[protocol], _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{})
	if err != nil {
		return err
	}
	defer m.conns[protocol].Close()

	quit := false
	m.conns[protocol].SetCloseHandler(func(code int, text string) error {
		quit = true
		return nil
	})

	d := wsInit{
		protocol,
		m.Name,
		m.Password,
	}

	err = m.conns[protocol].WriteJSON(d)
	if err != nil {
		return err
	}

	for !quit {
		var msg Message
		if m.conns[protocol].ReadJSON(&msg) != nil {
			continue
		}

		m.HSet.ServeMessage(m, msg)
	}

	return nil
}

func (m *Module) SendMessage(protocol, data string) {
	var msg Message
	msg.Protocol = protocol
	msg.Data = data
	msg.Module = m.Name

	m.conns[protocol].WriteJSON(msg)
}

func (m *Module) Close(protocol string) {
	m.conns[protocol].Close()
}

func (m *Module) CloseAll() {
	for _, v := range m.conns {
		if v != nil {
			v.Close()
		}
	}
}
