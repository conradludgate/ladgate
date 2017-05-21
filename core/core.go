package main

import (
	"encoding/json"
	"strconv"
	"strings"

	irc "github.com/fluffle/goirc/client"
)

func CoreOnInvite(conn *irc.Conn, line *irc.Line) {
	conn.Join(line.Text())
}

type Message struct {
	Error, Command, Channel, ToModule string
}

func CoreOnPrivMsg(conn *irc.Conn, line *irc.Line) {
	var req Message
	if err := json.Unmarshal([]byte(line.Text()), &req); err != nil {
		return
	}

	if req.ToModule != "core" {
		return
	}

	module := conn.StateTracker().GetChannel(line.Target()).Nicks[line.Nick]
	if module == nil {
		return
	}

	//args := strings.Split(req.Command, " ")
	command := irc.ParseLine(req.Command)

	switch command.Cmd {
	case "CONNECT":
		if !(module.Op || module.Owner) {
			resp := Message{"Only +o or above can run this command", "", "", line.Nick}
			msg, _ := json.Marshal(resp)

			conn.Privmsg(line.Target(), string(msg))
		}
		CoreConnect(conn, line, command, req)

	case "DISCONNECT":
		if !(module.Op || module.Owner) {
			resp := Message{"Only +o or above can run this command", "", "", line.Nick}
			msg, _ := json.Marshal(resp)

			conn.Privmsg(line.Target(), string(msg))
		}

		Connections[line.Target()].Quit("Goodbye.")
		Connections[line.Target()] = nil

	case "JOIN", "KICK", "NICK", "INVITE", "MODE":
		if !(module.Op || module.Owner || module.HalfOp) {
			resp := Message{"Only +h or above can run this command", "", "", line.Nick}
			msg, _ := json.Marshal(resp)

			conn.Privmsg(line.Target(), string(msg))
		}
		Connections[line.Target()].Raw(command.Raw)

	case "PRIVMSG", "ACTION":
		Connections[line.Target()].Raw(command.Raw)
	}
}

func CoreConnect(conn *irc.Conn, line, command *irc.Line, req Message) {
	if len(command.Args) == 0 {
		resp := Message{"CONNECT <server> [true|false]", "", req.Channel, line.Nick}
		msg, _ := json.Marshal(resp)

		conn.Privmsg(line.Target(), string(msg))
	}

	args := strings.Split(command.Text(), " ")

	var ssl bool
	if len(args) >= 1 {
		ssl, _ = strconv.ParseBool(args[1])
	}

	if err := NewConn(conn, line.Target(), args[0], ssl); err != nil {
		resp := Message{"Topic Error: " + err.Error(), "", "", ""}
		msg, _ := json.Marshal(resp)

		conn.Privmsg(line.Target(), string(msg))
		return
	}
}
