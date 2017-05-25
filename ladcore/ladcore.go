package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var adminHash []byte

func main() {
	if err := OpenModuleDB("ladcore.db"); err != nil {
		panic(err)
	}

	var key string
	key, adminHash = newKey(16)
	fmt.Println("Username: admin")
	fmt.Println("Password:", key)
	key = ""

	loadHttpServer(":2367", "", "")
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

func verifyModule(r *http.Request) (proto, module string, perms int) {
	user, pass, ok := r.BasicAuth()
	proto = r.FormValue("protocol")

	if !ok {
		return
	}

	module = user

	if module == "admin" {
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
	getHashStmt.QueryRow(module).Scan(&hash)

	if bcrypt.CompareHashAndPassword(hash, key) != nil {
		return
	}

	getPermsStmt.QueryRow(proto).Scan(&perms)

	return
}
