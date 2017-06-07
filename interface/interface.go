package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/conradludgate/ladgate/module"
	irc "github.com/fluffle/goirc/client"
)

func main() {
	ircmod = module.NewModule(module.LoadConfig("interface"))
	intface = module.NewModule(module.LoadConfig("admin"))

	module.HandleFunc("irc", HandleIRC)
	log.Fatal(ircmod.Listen("irc", nil))
}

var ircmod, intface *module.Module
var protocol string

func HandleIRC(m *module.Module, msg module.Message) {
	split := strings.SplitN(msg.Data, ": ", 2)
	if len(split) < 2 {
		return
	}

	server := split[0]
	data := split[1]

	command := irc.ParseLine(data)

	fmt.Println(command.Cmd, command.Text(), command.Target(), server)

	if command.Cmd == irc.PRIVMSG && command.Target() == "Oon" &&
		server == "zif" && command.Text()[0] == '@' {

		split = strings.SplitN(command.Text()[1:], ": ", 2)
		if len(split) != 2 {
			return
		}

		fmt.Println(split)

		switch split[0] {
		case "protocol":
			protocol = split[1]
			intface.CloseAll()
			go intface.Listen(protocol, module.HandlerFunc(HandleInterface))
		case "send":
			intface.SendMessage(protocol, split[1])
		}
	}
}

func HandleInterface(m *module.Module, msg module.Message) {
	ircmod.SendMessage("irc", "zif: PRIVMSG Oon "+msg.Module+": "+msg.Data)
}
