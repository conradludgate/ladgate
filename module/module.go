package module

import (
	"encoding/json"

	irc "github.com/fluffle/goirc/client"
)

type Module struct {
	conn     *irc.Conn
	handlers *hSet
}

type Message struct {
	Error, Command, Channel, ToModule string
}

type Line struct {
	Message Message

	Full, Line *irc.Line
}

func (l *Line) Copy() *Line {
	nl := Line{
		l.Message,
		l.Full.Copy(), l.Line.Copy(),
	}

	return &nl
}

func SimpleClient(nick string, args ...string) *Module {
	module := Client(irc.NewConfig(nick, args...))
	return module
}

func Client(cfg *irc.Config) *Module {
	conn := irc.Client(cfg)

	module := &Module{
		conn,
		handlerSet(),
	}

	conn.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		var msg Message
		if err := json.Unmarshal([]byte(line.Text()), &msg); err != nil {
			return
		}

		if msg.ToModule != "" && msg.ToModule != conn.Me().Nick {
			return
		}

		module.handlers.dispatch(module, &Line{msg, line, irc.ParseLine(msg.Command)})
	})

	conn.HandleFunc(irc.INVITE, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(line.Text())
	})

	return module
}

func (module *Module) Conn() *irc.Conn {
	return module.conn
}
