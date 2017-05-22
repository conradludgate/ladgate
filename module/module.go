package module

import (
	"encoding/json"

	irc "github.com/fluffle/goirc/client"
)

type Module struct {
	conn     *irc.Conn
	handlers *hSet
}

type Line struct {
	Error, Command, Channel, ToModule string

	Full, Line *irc.Line
}

func (l *Line) Copy() *Line {
	nl := Line{
		l.Error, l.Command, l.Channel, l.ToModule,
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
		var req Line
		if err := json.Unmarshal([]byte(line.Text()), &req); err != nil {
			return
		}

		if req.ToModule != conn.Me().Nick {
			return
		}

		req.Full = line
		req.Line = irc.ParseLine(req.Command)
		module.handlers.dispatch(module, &req)
	})

	conn.HandleFunc(irc.INVITE, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(line.Text())
	})

	return module
}

func (module *Module) Conn() *irc.Conn {
	return module.conn
}
