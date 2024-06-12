package db

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB_FILE = "./db/gTeam.db"
var CSV_FILE = "./db/employees.csv"

func InitDatabase() {

	db, err := sql.Open("sqlite", DB_FILE)
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
	log.Printf("Creating Tables ...")
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

}

func ImportEmployeesFromCSV(path string) {
	f, err := os.Open(path)

	if err != nil {
		log.Fatalf("open failed: %s", err)
	}

	r := csv.NewReader(f)
	// Read the header row.
	_, err = r.Read()
	if err != nil {
		log.Fatalf("missing header row(?): %s", err)
	}

	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		email := record[0]
		name := record[1]
		surname := record[2]
		team := record[3]

		sqlStmt, err := db.Prepare("insert into people(email, name, surname, team) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("INSERT prepare FAILED: %s", err)
		}

		_, err = sqlStmt.Exec(email, name, surname, team)
		if err != nil {
			log.Fatalf("INSERT has FAILED (%s): %s", email, err)
		}
	}
	log.Printf("Import of %s SUCCESS", path)
}
