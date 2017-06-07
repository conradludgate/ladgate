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

	proto, _, perms := getModulePermsWS(conn)

	if perms&GET == GET {
		conns[proto] = append(conns[proto], conn)
		//go subMessages(conn, proto, &quit)
	}

	if perms&SET == SET {
		for !quit {
			var m message
			if conn.ReadJSON(&m) != nil {
				return
			}

			if perms&2 == 2 {
				m.Perms = perms
				m.Time = time.Now()

				for i := 0; i < len(conns[proto]); i++ {
					v := conns[proto][i]
					if v == conn {
						continue
					}

					if v.WriteJSON(m) != nil {
						conns[proto] = append(conns[proto][:i], conns[proto][i+1:]...)
						i--
					}
				}
			}
		}
	}

	for !quit {
	}
}
