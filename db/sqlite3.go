package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() {
	db, err := sql.Open("sqlite3", "./gTeam.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT, email TEXT, team TEXT);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
