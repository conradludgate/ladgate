package main

import (
	"crypto/tls"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/conradludgate/ladgate/irc"
	"github.com/conradludgate/ladgate/module"

	"github.com/fluffle/goirc/client"
)

func main() {
	m := module.NewModule(module.LoadConfig("IRC"))

	patternHandler := module.NewPatternHandler()
	patternHandler.HandleFunc(irc.PatternCONNECT, HandleConnect).Perms(irc.PermCONNECT)
	patternHandler.HandleFunc(irc.PatternDISCONNECT, HandleDisconnect).Perms(irc.PermDISCONNECT)

	patternHandler.HandleFunc(irc.PatternJOIN, HandleRaw).Perms(irc.PermJOIN)
	patternHandler.HandleFunc(irc.PatternKICK, HandleRaw).Perms(irc.PermKICK)
	patternHandler.HandleFunc(irc.PatternNICK, HandleRaw).Perms(irc.PermNICK)
	patternHandler.HandleFunc(irc.PatternINVITE, HandleRaw).Perms(irc.PermINVITE)
	patternHandler.HandleFunc(irc.PatternMODE, HandleRaw).Perms(irc.PermMODE)

	patternHandler.HandleFunc(irc.PatternPRIVMSG, HandleRaw).Perms(irc.PermPRIVMSG)
	patternHandler.HandleFunc(irc.PatternACTION, HandleRaw).Perms(irc.PermACTION)

	err := m.Connect(irc.Protocol, irc.FinallyHandler{patternHandler})
	if err != nil {
		panic(err)
	}

	for {
	}
}

func HandleConnect(m *module.Module, message module.Message) {
	msg, _ := irc.ParseMessage(message.Data)

	if len(msg.IRC.Args) == 0 || Connections[msg.Server] != nil {
		return
	}

	cfg := client.NewConfig("ladgate", "ladgate", "Ladgate IRC Bridge")
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
	msg, _ := irc.ParseMessage(message.Data)

	conn, ok := Connections[msg.Server]
	if ok {
		conn.Quit("Goodbye.")
	}
}

func HandleRaw(m *module.Module, message module.Message) {
	msg, _ := irc.ParseMessage(message.Data)

	conn, ok := Connections[msg.Server]
	if ok {
		conn.Raw(msg.IRC.Raw)
	}
}
