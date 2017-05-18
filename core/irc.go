package main

import (
	"crypto/tls"
	"fmt"

	irc "github.com/fluffle/goirc/client"
)

type client struct {
	conn  *irc.Conn
	Admin string
}

type config struct {
	servers map[string]struct {
		SSL   bool
		Admin string
		Chans []string
	}
}

func main() {
	cfg := irc.NewConfig("ladgate", "Ladgate IRC Bot", "Developed by Oon")
	cfg.Version = "Ladgate 0.1"
	cfg.QuitMessage = "Piece out suckas!"

	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{ServerName: "irc.esper.net"}
	cfg.Server = "irc.esper.net:6697"

	c := irc.Client(cfg)

	c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		fmt.Println("Connected")
	})

	quit := make(chan bool)
	c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})

	c.HandleFunc(irc.INVITE, func(conn *irc.Conn, line *irc.Line) {
		fmt.Println(line.Public(), line.Text(), line.Args, line.Cmd, line.Ident)
		conn.Join(line.Text())
	})

	if err := c.ConnectTo("irc.esper.net"); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
		quit <- true
	}

	<-quit
}

func NewConn(servername, server string, ssl bool) {
	cfg := irc.NewConfig("ladgate", "Ladgate IRC Bot", "Developed by Oon")
	cfg.Version = "Ladgate 0.1"
	cfg.QuitMessage = "Piece out suckas!"

	cfg.SSL = ssl
	cfg.SSLConfig = &tls.Config{ServerName: servername}
	cfg.Server = server

	irc.Client(cfg)

}
