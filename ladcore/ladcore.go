package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var adminHash []byte

const (
	GET int = 1
	SET int = 2

	ADMIN      int = 1 << 31
	SUPERADMIN int = -1 << 31
)

func main() {
	if err := OpenModuleDB("ladcore.db"); err != nil {
		panic(err)
	}

	var key string
	key, adminHash = newKey(16)
	fmt.Println("Username: admin")
	fmt.Println("Password:", key)
	key = ""

	loadHttpServer(LoadConfig())
}

func newKey(cost int) (string, []byte) {
	key := make([]byte, 64)
	rand.Read(key)

	hash, err := bcrypt.GenerateFromPassword(key, cost)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(key), hash
}

func getModulePermsHTTP(r *http.Request) (string, string, int) {
	user, pass, ok := r.BasicAuth()
	proto := r.FormValue("proto")

	if !ok {
		return proto, user, 0
	}

	return proto, user, getModulePerms(proto, user, pass)
}

type wsInit struct {
	Protocol string `json:"protocol"`
	Module   string `json:"module"`
	Password string `json:"password"`
}

func getModulePermsWS(conn *websocket.Conn) (string, string, int) {
	var d wsInit

	if conn.ReadJSON(&d) != nil {
		return d.Protocol, d.Module, 0
	}

	return d.Protocol, d.Module, getModulePerms(d.Protocol, d.Module, d.Password)
}

func getModulePerms(proto, user, pass string) (perms int) {
	if user == "admin" {
		key, err := base64.StdEncoding.DecodeString(pass)
		if err != nil {
			return
		}

		if bcrypt.CompareHashAndPassword(adminHash, key) != nil {
			return
		}

		perms = -1 // All 1s (two's compliment). Admin has all potential permissions.
		return
	}

	key, err := base64.StdEncoding.DecodeString(pass)
	if err != nil {
		return
	}

	var hash []byte
	getHashStmt.QueryRow(user).Scan(&hash)

	if bcrypt.CompareHashAndPassword(hash, key) != nil {
		return
	}

	getPermsStmt.QueryRow(user, proto).Scan(&perms)

	return
}
