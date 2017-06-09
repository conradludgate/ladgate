package main

import (
	"crypto/tls"
	"flag"
	"strings"

	"github.com/conradludgate/ladgate/irc"
	"github.com/conradludgate/ladgate/module"
)

func main() {
	m := module.NewModule(module.LoadConfig("IRC"))

	patternHandler := module.NewPatternHandler()
	patternHandler.HandleFunc(irc.PatternCONNECT, HandleConnect).Perms(irc.PermCONNECT)
	patternHandler.HandleFunc(irc.PatternDISCONNECT, HandleDisconnect).Perms(irc.PermDISCONNECT)

	m.Connect(irc.Protocol, irc.FinallyHandler(patternHandler))
}

func HandleConnect(m *module.Module, message module.Message) {
	msg := irc.ParseMessage(message)

	if len(msg.IRC.Args) == 0 || Connections[msg.Server] != nil {
		return
	}

	cfg := irc.NewConfig("ladgate", "ladgate", "Ladgate IRC Bridge")
	cfg.Server = msg.IRC.Args[0]

	if len(msg.IRC.Args) >= 1 {
		flagset := flag.NewFlagSet("Connect", flag.ContinueOnError)

		ssl := flagset.BoolP("ssl", "s", false, "Enables SSL")
		sslVerify := flagset.Bool("ssl_verify", true, "Enables Verification of SSL")

		flagset.Parse(msg.IRC.Args[1:])

		cfg.SSL = *ssl

		if cfg.SSL {
			cfg.SSLConfig = &tls.Config{
				ServerName:         strings.Split(msg.Server, ":")[0],
				InsecureSkipVerify: !*sslVerify,
			}
		}
	}

	NewConn(m, msg.Server, cfg)
}

func HandleDisconnect(m *module.Module, message module.Message) {
	msg := irc.ParseMessage(message)

	conn, ok := Connections[msg.Server]
	if ok {
		conn.Quit("Goodbye.")
	}
}
