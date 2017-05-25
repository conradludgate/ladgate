package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type message struct {
	Data   string    `json:"data"`
	Module string    `json:"module"`
	Perms  int       `json:"perms"`
	Time   time.Time `json:"time"`
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

		messages[proto].addMessage(message{r.FormValue("data"), module, perms, time.Now()})

		io.WriteString(w, "Set message '"+r.FormValue("data")+"' on the protocol '"+proto+"'")

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func createNewModule(w http.ResponseWriter, r *http.Request) {
	_, _, perms := verifyModule(r)
	if perms&(-1<<31) == (-1 << 31) {
		key, hash := newKey(12)
		_, err := newModuleStmt.Exec(r.FormValue("module"), hash, r.FormValue("description"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		io.WriteString(w, "Username:\t"+r.FormValue("module")+"\nPassword:\t"+key)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func givePerms(w http.ResponseWriter, r *http.Request) {
	proto, _, perms := verifyModule(r)
	if perms&(-1<<31) == (-1 << 31) {
		p, err := strconv.Atoi(r.FormValue("perms"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		if p < 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Only the admin (you) can have a negative permission")
			return
		}

		if _, err := replacePermsStmt.Exec(r.FormValue("module"), proto, p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		io.WriteString(w, r.FormValue("module")+" has been granted the permissions "+r.FormValue("perms")+" for the '"+proto+"' protocol")

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
