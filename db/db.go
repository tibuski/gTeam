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

func InitDatabase(DB_FILE string) {

	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Printf("Failed to open/create file %s with error %s", DB_FILE, err)
	}

	defer db.Close()

	sqlStmt := `
				CREATE TABLE IF NOT EXISTS people (employeeNumber INTEGER PRIMARY KEY, email TEXT, name TEXT, surname TEXT, team TEXT, UNIQUE (email));
				CREATE TABLE IF NOT EXISTS eventTypes (code INTEGER PRIMARY KEY, description TEXT, descriptionLong TEXT, picture TEXT);
				CREATE TABLE IF NOT EXISTS dutyTypes (code INTEGER PRIMARY KEY, description TEXT, descriptionLong TEXT, picture TEXT);
				CREATE TABLE IF NOT EXISTS eventTable (employeeNumber TEXT, code INTEGER, day TEXT, month TEXT, year TEXT, timestamp TEXT, who TEXT, deleted INT);
				CREATE TABLE IF NOT EXISTS dutyTable (employeeNumber TEXT, code INTEGER, day TEXT, month TEXT, year TEXT, timestamp TEXT, who TEXT, deleted INT);
				`
	_, err = db.Exec(sqlStmt)
	log.Printf("Creating Tables ...")
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func ImportEmployeesFromCSV(DB_FILE string, csvPath string) {
	f, err := os.Open(csvPath)
	if err != nil {
		log.Printf("Failed to open file %s with error : %s", csvPath, err)
	}

	r := csv.NewReader(f)
	// Read the header row.
	_, err = r.Read()
	if err != nil {
		log.Printf("missing header row(?): %s", err)
	}

	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Printf("Failed to open/create file %s with error %s", DB_FILE, err)
	}

	defer db.Close()

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		employeeNumber := record[0]
		email := record[1]
		name := record[2]
		surname := record[3]
		team := record[4]

		sqlStmt, err := db.Prepare("INSERT INTO people(employeeNumber, email, name, surname, team) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("INSERT prepare FAILED: %s", err)
		}

		_, err = sqlStmt.Exec(employeeNumber, email, name, surname, team)
		if err != nil {
			log.Printf("INSERT has FAILED : %s", err)
		}
	}
	log.Printf("Import of %s DONE", f.Name())
}

func ImportTypesFromCSV(DB_FILE string, csvPath string, table string) {

	var switchStmt string
	switch table {
	case "eventTypes":
		switchStmt = "INSERT INTO eventTypes(code, description, descriptionLong, picture) VALUES(?, ?, ?, ?)"

	case "dutyTypes":
		switchStmt = "INSERT INTO dutyTypes(code, description, descriptionLong, picture) VALUES(?, ?, ?, ?)"

	}

	f, err := os.Open(csvPath)
	if err != nil {
		log.Printf("Failed to open file %s with error : %s", csvPath, err)
	}

	r := csv.NewReader(f)
	// Read the header row.
	_, err = r.Read()
	if err != nil {
		log.Printf("missing header row(?): %s", err)
	}

	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Printf("Failed to open/create file %s with error %s", DB_FILE, err)
	}

	defer db.Close()

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		code := record[0]
		description := record[1]
		descriptionLong := record[2]
		picture := record[3]
		sqlStmt, err := db.Prepare(switchStmt)
		if err != nil {
			log.Printf("INSERT prepare FAILED: %s", err)
		}

		_, err = sqlStmt.Exec(code, description, descriptionLong, picture)
		if err != nil {
			log.Printf("INSERT has FAILED : %s", err)
		}
	}
	log.Printf("Import of %s DONE", f.Name())
}

func ImportTablesFromCSV(DB_FILE string, csvPath string, table string) {

	var switchStmt string
	switch table {
	case "eventTable":
		switchStmt = "INSERT INTO eventTable(employeeNumber, code, day, month, year, timestamp, who, deleted) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	case "dutyTable":
		switchStmt = "INSERT INTO dutyTable(employeeNumber, code, day, month, year, timestamp, who, deleted) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	}

	f, err := os.Open(csvPath)
	if err != nil {
		log.Printf("Failed to open file %s with error : %s", csvPath, err)
	}

	r := csv.NewReader(f)
	// Read the header row.
	_, err = r.Read()
	if err != nil {
		log.Printf("missing header row(?): %s", err)
	}

	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Printf("Failed to open/create file %s with error %s", DB_FILE, err)
	}

	// Optimize SQLite for faster writes
	_, err = db.Exec("PRAGMA synchronous = OFF; PRAGMA journal_mode = MEMORY;")
	if err != nil {
		log.Printf("Failed to optimize SQLite: %s", err)
		return
	}

	defer db.Close()

	var i = 1

	for {

		record, err := r.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		employeeNumber := record[0]
		code := record[1]
		day := record[2]
		month := record[3]
		year := record[4]
		timestamp := record[5]
		who := record[6]
		deleted := record[7]
		sqlStmt, err := db.Prepare(switchStmt)
		if err != nil {
			log.Printf("INSERT prepare FAILED: %s", err)
		}

		_, err = sqlStmt.Exec(employeeNumber, code, day, month, year, timestamp, who, deleted)
		if err != nil {
			log.Printf("INSERT has FAILED : %s", err)
		}
		i++
		if i%1000 == 0 {
			log.Printf("Imported %d records in %s", i, table)
		}
	}
	log.Printf("Import of %s DONE", f.Name())
}

func SelectFromTable(DB_FILE string, sqlStmt string) ([]map[string]interface{}, error) {

	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Printf("Failed to open file %s with error %s", DB_FILE, err)
	}
	defer db.Close()

	// Optimize SQLite for faster writes
	_, err = db.Exec("PRAGMA synchronous = OFF; PRAGMA journal_mode = MEMORY;")
	if err != nil {
		log.Printf("Failed to optimize SQLite: %s", err)
		return nil, err
	}

	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("Failed to execute query %q: %s", sqlStmt, err)
		return nil, err
	}
	defer rows.Close()

	log.Printf("Executed query: %s", sqlStmt)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Create a slice of maps to hold the query results
	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interfaces to hold each column value
		values := make([]interface{}, len(columns))
		// Create a slice of pointers to each value
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		// Create a map to represent the row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
