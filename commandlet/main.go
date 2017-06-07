package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/conradludgate/ladgate/module"
	irc "github.com/fluffle/goirc/client"
)

func main() {
	m := module.NewModule(module.LoadConfig("cmdlet"))

	module.HandleFunc("irc", HandleIRC)

	// go log.Fatal(m.Listen("messenger", nil))
	// go log.Fatal(m.Listen("discord", nil))
	log.Fatal(m.Listen("irc", nil))
}

func HandleIRC(m *module.Module, msg module.Message) {
	split := strings.SplitN(msg.Data, ": ", 2)
	if len(split) < 2 {
		return
	}

	server := split[0]
	data := split[1]

	command := irc.ParseLine(data)
	if command.Cmd != irc.PRIVMSG {
		return
	}

	text := strings.Trim(command.Text(), " ")
	if len(text) < 2 {
		return
	}
	if text[0] != '!' {
		return
	}

	split = strings.SplitN(text[1:], " ", 2)

	input := ""
	if len(split) == 2 {
		input = split[1]
	}

	r := Commandlet(split[0], input, command.Nick)
	if r != "" {
		m.SendMessage("irc", fmt.Sprintf("%s: PRIVMSG %s :%s", server, command.Target(), r))
	}

}

func Commandlet(command, input, user string) string {
	switch command {
	case "kys":
		return "No! We love you, " + user + "!"
	case "?":
		return "Help? You don't need help! /wjh make some docs/"
	case "!":
		return "Place holder for adding commandlets and shit"
	}

	return ""
}
