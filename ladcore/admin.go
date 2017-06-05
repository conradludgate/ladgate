package main

import (
	"io"
	"net/http"
	"strconv"
)

func createNewModule(w http.ResponseWriter, r *http.Request) {
	_, _, perms := getModulePermsHTTP(r)
	if perms&ADMIN == ADMIN {
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
	proto, _, perms := getModulePermsHTTP(r)
	if perms&ADMIN == ADMIN {
		p, err := strconv.Atoi(r.FormValue("perms"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		if p&SUPERADMIN == SUPERADMIN {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Only the superadmin can have a negative permission")
			return
		}

		if p&ADMIN == ADMIN && perms&SUPERADMIN != SUPERADMIN {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Only the superadmin can make more admins")
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
