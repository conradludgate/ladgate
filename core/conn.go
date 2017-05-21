package main

import (
	"crypto/tls"
	"encoding/json"
	"strings"

	irc "github.com/fluffle/goirc/client"
	"github.com/spf13/viper"
)

var Connections map[string]*irc.Conn

func NewConn(core *irc.Conn, channel, server string, ssl bool) error {
	cfg := irc.NewConfig(viper.GetString("Nick"),
		viper.GetString("Ident"),
		viper.GetString("Name"))

	cfg.Server = server
	cfg.SSL = ssl

	if ssl {
		cfg.SSLConfig = &tls.Config{ServerName: strings.Split(server, ":")[0]}
	}

	c := irc.Client(cfg)

	Connections[channel] = c

	onEvent := func(conn *irc.Conn, line *irc.Line) {
		resp := Message{"", line.Raw, line.Target(), ""}
		msg, _ := json.Marshal(resp)

		core.Privmsg(channel, string(msg))
	}

	c.HandleFunc(irc.REGISTER, onEvent)
	//c.HandleFunc(irc.CONNECTED, onEvent)
	c.HandleFunc(irc.DISCONNECTED, onEvent)
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

	return c.Connect()
}
