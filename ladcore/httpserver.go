package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type message struct {
	Protocol string    `json:"protocol"`
	Data     string    `json:"data"`
	Module   string    `json:"module"`
	Perms    int       `json:"perms"`
	Time     time.Time `json:"time"`
}

const storedMessages int = 50

type messageQueue struct {
	messages [storedMessages]message

	FrontPointer, RearPointer int
}

func (q *messageQueue) addMessage(m message) {
	q.messages[q.RearPointer] = m

	q.RearPointer += 1
	q.RearPointer %= storedMessages

	if q.RearPointer == q.FrontPointer {
		q.FrontPointer += 1
		q.FrontPointer %= storedMessages
	}
}

var messages map[string]*messageQueue

func loadHttpServer(addr, certFile, keyFile string) {
	messages = make(map[string]*messageQueue)

	http.HandleFunc("/get", getMessages)
	http.HandleFunc("/set", setMessage)

	http.HandleFunc("/new", createNewModule)
	http.HandleFunc("/perm", givePerms)

	//http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	http.ListenAndServe(addr, nil)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	proto, _, perms := verifyModule(r)
	if perms&1 == 1 {
		q := messages[proto]
		if q == nil {
			q = &messageQueue{}
		}

		var m []message
		if q.FrontPointer <= q.RearPointer {
			m = q.messages[q.FrontPointer:q.RearPointer]
		} else {
			m = append(q.messages[q.FrontPointer:], q.messages[:q.RearPointer]...)
		}

		b, _ := json.MarshalIndent(m, "", "\t")
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func setMessage(w http.ResponseWriter, r *http.Request) {
	proto, module, perms := verifyModule(r)
	if perms&2 == 2 {
		if messages[proto] == nil {
			messages[proto] = &messageQueue{}
		}

		messages[proto].addMessage(message{proto, r.FormValue("data"), module, perms, time.Now()})

		io.WriteString(w, "Set message '"+r.FormValue("data")+"' on the protocol '"+proto+"'")

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func getPerms(w http.ResponseWriter, r *http.Request) {

}
