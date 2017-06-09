package main

import (
	"io"
	"strings"

	"github.com/conradludgate/ladgate/irc"
	"github.com/conradludgate/ladgate/module"
)

func main() {
	ircModule := module.NewModule(module.LoadConfig("cmdlet"))

	ircHandler := module.NewPatternHandler()
	ircHandler.HandleFunc(irc.PatternPRIVMSG)

	module.HandleFunc(module.Pattern(MatchCommandlet), HandleIRC).Perms(irc.PermIRCBridge)

	// go log.Fatal(m.Listen("messenger", nil))
	// go log.Fatal(m.Listen("discord", nil))
	ircModule.Connect(irc.Protocol, irc.FinallyHandler{ircHandler})

	for {
	}
}

func MatchCommandlet(message string) (match bool) {
	if irc.PatternPRIVMSG.Match(message) {
		match = true

		msg, _ := irc.ParseMessage(message)
		text := strings.Trim(msg.Text(), " ")
		if len(text) < 2 {
			return false
		}
		if text[0] != '!' {
			match = false
		}
	}
	return
}

func HandleIRC(m *module.Module, message module.Message) {
	msg := irc.ParseMessage(message)

	text := strings.Trim(msg.Text(), " ")
	split = strings.SplitN(text[1:], " ", 2)

	input := ""
	if len(split) == 2 {
		input = split[1]
	}

	r := Commandlet(split[0], input, command.Nick)
	if r != "" {
		io.WriteString(m, msg.Server+": PRIVMSG "+msg.IRC.Target()+" :"+r)
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
