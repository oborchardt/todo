package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB = nil

func getDb() *sql.DB {
	if db == nil {
		conn, err := sql.Open("sqlite3", "todo.db")
		if err != nil {
			log.Fatal(err)
		}
		db = conn
	}
	return db
}
