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

type People struct {
	EmployeeNumber int
	Email          string
	Name           string
	Surname        string
	Team           string
}

type PeopleUnion struct {
	EmployeeNumber int
	Email          string
	Name           string
	Surname        string
	Team           string
	DutyCode       string
	EventCode      string
	Day            int
	Month          int
	Year           int
	Timestamp      string
	Who            string
	Deleted        int
}

func OpenDatabase(DB_FILE string) *sql.DB {

	database, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		log.Printf("Failed to open/create file %s with error %s", DB_FILE, err)
	}

	// Optimize SQLite for faster writes
	_, err = database.Exec("PRAGMA synchronous = OFF; PRAGMA journal_mode = MEMORY;")
	if err != nil {
		log.Printf("Failed to optimize SQLite: %s", err)
	}

	return database
}

func InitDatabase(database *sql.DB) {

	sqlStmt := `
                CREATE TABLE IF NOT EXISTS people (employeeNumber INTEGER PRIMARY KEY, email TEXT, name TEXT, surname TEXT, team TEXT, UNIQUE (email));
                CREATE TABLE IF NOT EXISTS eventTypes (code INTEGER PRIMARY KEY, description TEXT, descriptionLong TEXT, picture TEXT);
                CREATE TABLE IF NOT EXISTS dutyTypes (code INTEGER PRIMARY KEY, description TEXT, descriptionLong TEXT, picture TEXT);
                CREATE TABLE IF NOT EXISTS eventTable (employeeNumber TEXT, code INTEGER, day TEXT, month TEXT, year TEXT, timestamp TEXT, who TEXT, deleted INT);
                CREATE TABLE IF NOT EXISTS dutyTable (employeeNumber TEXT, code INTEGER, day TEXT, month TEXT, year TEXT, timestamp TEXT, who TEXT, deleted INT);
                `
	_, err := database.Exec(sqlStmt)
	log.Printf("Creating Tables ...")
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func ImportEmployeesFromCSV(database *sql.DB, csvPath string) {
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

		sqlStmt, err := database.Prepare("INSERT INTO people(employeeNumber, email, name, surname, team) VALUES(?, ?, ?, ?, ?)")
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

func ImportTypesFromCSV(database *sql.DB, csvPath string, table string) {

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

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		code := record[0]
		description := record[1]
		descriptionLong := record[2]
		picture := record[3]
		sqlStmt, err := database.Prepare(switchStmt)
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

func ImportTablesFromCSV(database *sql.DB, csvPath string, table string) {

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
		sqlStmt, err := database.Prepare(switchStmt)
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

func SelectFromPeople(database *sql.DB, filter string) ([]People, error) {

	var Peoples []People

	sqlStmt := `SELECT * FROM people WHERE employeenumber LIKE ?`

	rows, err := database.Query(sqlStmt, filter)
	if err != nil {
		log.Printf("Failed to execute query %q: %s", sqlStmt, err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var p People
		err := rows.Scan(&p.EmployeeNumber, &p.Email, &p.Name, &p.Surname, &p.Team)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		Peoples = append(Peoples, p)
	}

	// Check for any errors encountered during the iteration of rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
		return nil, err
	}

	// Return the slice of People and a nil error indicating success
	return Peoples, nil
}

func SelectPeopleAllEvents(database *sql.DB, month, year int) ([]PeopleUnion, error) {

	var PeoplesAllEvents []PeopleUnion

	sqlStmt := `
                SELECT 
                    p.employeeNumber,
                    p.email,
                    p.name,
                    p.surname,
                    p.team,
                    d.code AS dutyCode,
                    0 AS eventCode,
                    d.day, d.month, d.year,
                    d.timestamp,
                    d.who,
                    d.deleted

                FROM 
                    people p
                JOIN 
                    dutyTable d ON p.employeeNumber = d.employeeNumber
                WHERE 
                    d.month = ? AND d.year = ?

                UNION ALL

                SELECT 
                    p.employeeNumber,
                    p.email,
                    p.name,
                    p.surname,
                    p.team,
                    0 AS dutyCode,
                    e.code AS eventCode,
                    e.day, e.month, e.year,
                    e.timestamp,
                    e.who,
                    e.deleted
                FROM 
                    people p
                JOIN 
                    eventTable e ON p.employeeNumber = e.employeeNumber
                WHERE 
                    e.month = ? AND e.year = ?
                `

	rows, err := database.Query(sqlStmt, month, year, month, year)
	if err != nil {
		log.Printf("Failed to execute query %q: %s", sqlStmt, err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var p PeopleUnion
		err := rows.Scan(&p.EmployeeNumber, &p.Email, &p.Name, &p.Surname, &p.Team, &p.DutyCode, &p.EventCode, &p.Day, &p.Month, &p.Year, &p.Timestamp, &p.Who, &p.Deleted)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		PeoplesAllEvents = append(PeoplesAllEvents, p)
	}

	// Check for any errors encountered during the iteration of rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
		return nil, err
	}

	// Return the slice of People and a nil error indicating success
	return PeoplesAllEvents, nil
}
