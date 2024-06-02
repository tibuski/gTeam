package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() {
	db, err := sql.Open("sqlite3", "./gTeam.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
				CREATE TABLE IF NOT EXISTS people (
					email TEXT PRIMARY KEY, 
					name TEXT, 
					surname TEXT, 
					team TEXT, 
					UNIQUE (email)
					);
				`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func importFromCSV(f string) {

}
