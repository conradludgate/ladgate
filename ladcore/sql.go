package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenModuleDB(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	db.Exec(dbInit)
	getHashStmt, _ = db.Prepare(getHash)
	getPermsStmt, _ = db.Prepare(getPerms)
	replacePermsStmt, _ = db.Prepare(replacePerms)
	newModuleStmt, _ = db.Prepare(newModule)
	updateHashStmt, _ = db.Prepare(updateHash)

	return nil
}

const (
	dbInit = `CREATE TABLE OR IGNORE
	module (
		name STRING PRIMARY KEY,
		hash BLOB,
		desc STRING
	);

	CREATE TABLE OR IGNORE
	perms (
		name STRING,
		protocol STRING,
		perms INTEGER,

		PRIMARY KEY(name, protocol)
	);`

	getHash = `SELECT (hash)
	FROM module
	WHERE name=?`

	getPerms = `SELECT perms
	FROM perms
	WHERE name=?
	AND protocol=?`

	replacePerms = `REPLACE INTO
	perms(name, protocol, perms)
	VALUES(?,?,?);`

	newModule = `INSERT INTO module 
	(name, hash, desc) 
	VALUES(?,?,?);`

	updateHash = `UPDATE module
	SET hash = ?
	WHERE name = ?;`
)

var getHashStmt *sql.Stmt
var getPermsStmt *sql.Stmt
var replacePermsStmt *sql.Stmt
var newModuleStmt *sql.Stmt
var updateHashStmt *sql.Stmt
