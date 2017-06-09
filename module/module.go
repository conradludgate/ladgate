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

// Implements WriteCloser
type Module struct {
	Server         string
	Name, Password string

	Handler Handler

	conn *websocket.Conn
}

func NewModule(server, name, password string) *Module {
	var m Module
	m.Server = server
	m.Name = name
	m.Password = password

	return &m
}

func (m *Module) Connect(protocol string, h Handler) error {
	if h == nil {
		m.Handler = DefaultPatternHandler
	} else {
		m.Handler = h
	}

	u, err := url.Parse(m.Server)
	if err != nil {
		return err
	}

	u.Path = "ws"

	m.conn, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{})
	if err != nil {
		return err
	}

	var quit bool
	m.conn.SetCloseHandler(func(code int, text string) error {
		quit = true
		return nil
	})

	d := `{"protocol": "` + protocol + `", "module": "` + m.Name +
		`", "password": "` + m.Password + `"}`

	err = m.conn.WriteMessage(websocket.TextMessage, []byte(d))
	if err != nil {
		return err
	}

	go func() {
		for !quit {
			var msg Message
			if m.conn.ReadJSON(&msg) != nil {
				continue
			}

			m.Handler.ServeMessage(m, msg)
		}
	}()

	return nil
}

// Implements io.Writer
func (m *Module) Write(p []byte) (n int, err error) {
	err = m.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Implements io.Closer,
func (m *Module) Close() error {
	return m.conn.Close()
}
