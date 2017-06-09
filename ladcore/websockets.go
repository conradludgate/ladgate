package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websockets(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	quit := false
	conn.SetCloseHandler(func(code int, text string) error {
		quit = true
		return nil
	})

	proto, module, perms := getModulePermsWS(conn)

	if perms&3 == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, ok := conns[proto]
	if !ok {
		conns[proto] = make(map[string]*websocket.Conn)
	}

	if perms&GET == GET {
		conns[proto][module] = conn
	}

	if perms&SET == SET {
		for !quit {
			var m message
			if conn.ReadJSON(&m) != nil {
				return
			}

			m.Perms = perms
			m.Time = time.Now()

			for k, v := range conns[proto] {
				if k == module {
					continue
				}

				if v.WriteJSON(m) != nil {
					delete(conns[proto], k)
				}
			}
		}
	}

	for !quit {
	}
}
