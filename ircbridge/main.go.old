package main

import (
	"crypto/tls"
	"log"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/conradludgate/ladgate/module"
	irc "github.com/fluffle/goirc/client"
)

const (
	PRIVMSG int = 4
	ACTION
	NAMES

	JOIN int = 8
	KICK
	NICK
	INVITE
	MODE

	CONNECT int = 16
	DISCONNECT
)

func main() {
	m := module.NewModule(module.LoadConfig("IRC"))

	module.HandleFunc("irc", HandleMessage)

	log.Fatal(m.Listen("irc", nil))
}

func HandleMessage(m *module.Module, msg module.Message) {
	split := strings.SplitN(msg.Data, ": ", 2)
	if len(split) < 2 {
		return
	}

	server := split[0]
	data := split[1]

	conn, ok := Connections[server]

	command := irc.ParseLine(data)

	switch command.Cmd {
	case "CONNECT":
		if msg.Perms&CONNECT == CONNECT && !ok {
			Connect(m, server, command)
		}
	case "DISCONNECT":
		if msg.Perms&DISCONNECT == DISCONNECT && ok {
			conn.Quit("Goodbye.")
		}
	case irc.JOIN, irc.KICK, irc.NICK, irc.INVITE, irc.MODE:
		if msg.Perms&JOIN == JOIN && ok {
			conn.Raw(command.Raw)
		}
	case irc.PRIVMSG, irc.ACTION:
		if msg.Perms&PRIVMSG == PRIVMSG && ok {
			conn.Raw(command.Raw)
		}
	case "NAMES":
		if msg.Perms&NAMES == NAMES && ok {
			Names(m, command, conn)
		}
	}
}

func Connect(m *module.Module, server string, command *irc.Line) {
	if len(command.Args) == 0 || Connections[server] != nil {
		return
	}

	cfg := irc.NewConfig("ladgate", "ladgate", "Ladgate IRC Bridge")
	cfg.Server = command.Args[0]

	if len(command.Args) >= 1 {
		flagset := flag.NewFlagSet("Connect", flag.ContinueOnError)

		ssl := flagset.BoolP("ssl", "s", false, "Enables SSL")
		sslVerify := flagset.Bool("ssl_verify", true, "Enables Verification of SSL")

		flagset.Parse(command.Args[1:])

		cfg.SSL = *ssl

		if cfg.SSL {
			cfg.SSLConfig = &tls.Config{
				ServerName:         strings.Split(server, ":")[0],
				InsecureSkipVerify: !*sslVerify,
			}
		}
	}

	NewConn(m, server, cfg)
}

func Names(m *module.Module, command *irc.Line, conn *irc.Conn) {
	//conn.StateTracker().GetChannel(command.Text()).Nicks
}
