package main

import (
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

var conns = make(map[string][]*websocket.Conn)

func loadHttpServer(addr, certFile, keyFile string) {
	http.HandleFunc("/set", BasicAuth(setMessage))

	http.HandleFunc("/new", BasicAuth(createNewModule))
	http.HandleFunc("/perm", BasicAuth(givePerms))

	http.HandleFunc("/ws", websockets)

	//http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	http.ListenAndServe(addr, nil)
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

		for i := 0; i < len(conns[proto]); i++ {
			v := conns[proto][i]
			if v.WriteJSON(m) != nil {
				conns[proto] = append(conns[proto][:i], conns[proto][i+1:]...)
				i--
			}
		}

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
