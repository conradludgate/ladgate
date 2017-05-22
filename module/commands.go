package module

import "encoding/json"

func (module *Module) Raw(line *Line, rawstring string) {
	msg := Message{"", rawstring, line.Message.Channel, line.Full.Nick}
	b, _ := json.Marshal(msg)

	module.conn.Privmsg(line.Full.Target(), string(b))
}
