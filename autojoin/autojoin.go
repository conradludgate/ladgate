package main

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	irc "github.com/fluffle/goirc/client"

	"github.com/conradludgate/ladgate/module"
)

func main() {
	cfg := irc.NewConfig("autojoin", "autojoin", "Autojoin Ladgate Module")

	cfg.Server = "ladgate.conradludgate.com:6697"
	cfg.SSL = true

	if cfg.SSL {
		cfg.SSLConfig = &tls.Config{ServerName: strings.Split(cfg.Server, ":")[0]}
	}

	m := module.Client(cfg)

	m.HandleFunc(irc.DISCONNECTED, func(module *module.Module, line *module.Line) {
		time.Sleep(time.Second * 10)
		m.Conn().Connect()
	})

	m.HandleFunc(irc.INVITE, OnInvite)

	quit := make(chan bool)
	if err := m.Conn().Connect(); err != nil {
		fmt.Println(err.Error())
		quit <- true
	}

	<-quit
}

func OnInvite(module *module.Module, line *module.Line) {
	module.Raw(line, irc.JOIN+" "+line.Line.Text())
}
