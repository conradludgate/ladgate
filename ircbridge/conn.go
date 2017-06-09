package main

import (
	"io"

	"github.com/conradludgate/ladgate/module"
	irc "github.com/fluffle/goirc/client"
)

var Connections = make(map[string]*irc.Conn)

func NewConn(m *module.Module, server string, cfg *irc.Config) {
	c := irc.Client(cfg)

	Connections[server] = c

	onEvent := func(conn *irc.Conn, line *irc.Line) {
		io.WriteString(m, server+": "+line.Raw)
	}

	c.EnableStateTracking()

	c.HandleFunc(irc.REGISTER, onEvent)
	//c.HandleFunc(irc.CONNECTED, onEvent)
	//c.HandleFunc(irc.DISCONNECTED, onEvent)
	c.HandleFunc(irc.ACTION, onEvent)
	c.HandleFunc(irc.AWAY, onEvent)
	c.HandleFunc(irc.CAP, onEvent)
	c.HandleFunc(irc.CTCP, onEvent)
	c.HandleFunc(irc.CTCPREPLY, onEvent)
	c.HandleFunc(irc.INVITE, onEvent)
	c.HandleFunc(irc.JOIN, onEvent)
	c.HandleFunc(irc.KICK, onEvent)
	c.HandleFunc(irc.MODE, onEvent)
	c.HandleFunc(irc.NICK, onEvent)
	//c.HandleFunc(irc.NOTICE, onEvent)
	//c.HandleFunc(irc.OPER, onEvent)
	c.HandleFunc(irc.PART, onEvent)
	//c.HandleFunc(irc.PASS, onEvent)
	//c.HandleFunc(irc.PING, onEvent)
	//c.HandleFunc(irc.PONG, onEvent)
	c.HandleFunc(irc.PRIVMSG, onEvent)
	c.HandleFunc(irc.QUIT, onEvent)
	c.HandleFunc(irc.TOPIC, onEvent)
	//c.HandleFunc(irc.USER, onEvent)
	//c.HandleFunc(irc.VERSION, onEvent)
	//c.HandleFunc(irc.VHOST, onEvent)
	c.HandleFunc(irc.WHO, onEvent)
	c.HandleFunc(irc.WHOIS, onEvent)

	onQuit := func(conn *irc.Conn, line *irc.Line) {
		conn.Close()
		Connections[server] = nil
	}

	c.HandleFunc(irc.DISCONNECTED, onQuit)
	c.HandleFunc(irc.QUIT, onQuit)

	c.HandleFunc(irc.INVITE, func(conn *irc.Conn, line *irc.Line) {
		conn.Join(line.Text())
	})

	c.Connect()
}
