package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type message struct {
	Protocol string    `json:"protocol"`
	Data     string    `json:"data"`
	Module   string    `json:"module"`
	Perms    int       `json:"perms"`
	Time     time.Time `json:"time"`
}

var conns = make(map[string]map[string]*websocket.Conn)

func loadHttpServer(addr, certFile, keyFile string) {
	http.HandleFunc("/set", BasicAuth(setMessage))

	http.HandleFunc("/new", BasicAuth(createNewModule))
	http.HandleFunc("/perm", BasicAuth(givePerms))

	http.HandleFunc("/ws", websockets)

	log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, nil))
}

func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, ok := r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Ladcore: Authentication Required"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func setMessage(w http.ResponseWriter, r *http.Request) {
	proto, module, perms := getModulePermsHTTP(r)
	if perms&SET == SET {

		m := message{proto, r.FormValue("data"), module, perms, time.Now()}

		_, ok := conns[proto]
		if !ok {
			return
		}

		for k, v := range conns[proto] {
			if k == module {
				continue
			}

			if v.WriteJSON(m) != nil {
				delete(conns[proto], k)
			}
		}

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
