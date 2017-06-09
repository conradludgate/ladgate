package irc

import (
	"errors"
	"strings"

	"github.com/conradludgate/ladgate/module"
	irc "github.com/fluffle/goirc/client"
)

var (
	PatternPRIVMSG module.PatternFunc = module.PatternFunc(matchPrivmsg)
	PatternACTION  module.PatternFunc = module.PatternFunc(matchAction)

	PatternJOIN   module.PatternFunc = module.PatternFunc(matchJoin)
	PatternKICK   module.PatternFunc = module.PatternFunc(matchKick)
	PatternNICK   module.PatternFunc = module.PatternFunc(matchNick)
	PatternINVITE module.PatternFunc = module.PatternFunc(matchInvite)
	PatternMODE   module.PatternFunc = module.PatternFunc(matchMode)

	PatternCONNECT    module.PatternFunc = module.PatternFunc(matchConnect)
	PatternDISCONNECT module.PatternFunc = module.PatternFunc(matchDisconnect)
)

type IRCMessage struct {
	Server string
	IRC    *irc.Line
}

var Messages = make(map[string]IRCMessage)

var ErrBadMessage = errors.New("Bad message format")

func ParseMessage(message string) (IRCMessage, error) {
	m, ok := Messages[message]
	if ok {
		return m, nil
	}

	split := strings.SplitN(message, ": ", 2)
	if len(split) < 2 {
		return m, ErrBadMessage
	}

	m.Server = split[0]
	m.IRC = irc.ParseLine(split[1])

	Messages[message] = m
	return m, nil
}

type FinallyHandler struct {
	Handle module.Handler
}

func Finally(message string) {
	delete(Messages, message)
}

func (fh *FinallyHandler) ServeMessage(m *module.Module, message module.Message) {
	fh.Handle.ServeMessage(m, message)
	Finally(message.Data)
}

func matchMessage(message, command string) bool {
	ircMessage, err := ParseMessage(message)
	if err != nil {
		return false
	}
	return ircMessage.IRC.Cmd == command
}

func matchPrivmsg(message string) bool {
	return matchMessage(message, irc.PRIVMSG)
}

func matchAction(message string) bool {
	return matchMessage(message, irc.ACTION)
}

func matchJoin(message string) bool {
	return matchMessage(message, irc.JOIN)
}

func matchKick(message string) bool {
	return matchMessage(message, irc.KICK)
}

func matchNick(message string) bool {
	return matchMessage(message, irc.NICK)
}

func matchInvite(message string) bool {
	return matchMessage(message, irc.INVITE)
}

func matchMode(message string) bool {
	return matchMessage(message, irc.MODE)
}

func matchConnect(message string) bool {
	return matchMessage(message, "CONNECT")
}

func matchDisconnect(message string) bool {
	return matchMessage(message, "DISCONNECT")
}
