package main

import (
	"fmt"
	"log"

	"github.com/tibuski/gTeam/db"
)

type Employee struct {
	employeeNumber int
	email          string
	name           string
	surname        string
	team           string
}

type Calendar struct {
	employeeNumber int
	days           []string
}

var Employees = []Employee{}
var TeamCalendar = []Calendar{}

const DB_FILE = "./db/gTeam.db"
const EMPLOYEES_CSV = "./db/people.csv"
const EVENTTYPES_CSV = "./db/eventTypes.csv"
const EVENTTABLE_CSV = "./db/eventTable.csv"
const DUTYTYPES_CSV = "./db/dutyTypes.csv"
const DUTYTABLE_CSV = "./db/dutyTable.csv"

// func daysOfTheMonth(month int, year int) []string {
// 	var days []string
// 	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
// 	end := start.AddDate(0, 1, -1)
// 	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
// 		if t.Weekday() == 6 || t.Weekday() == 0 {
// 			days = append(days, "W")
// 		} else {
// 			days = append(days, "0")
// 		}
// 	}
// 	return days
// }

// func createCalendar(month, year int) []Calendar {
// 	d := daysOfTheMonth(month, year)
// 	for _, n := range Employees {
// 		TeamCalendar = append(TeamCalendar, Calendar{employeeNumber: n.ame, Days: d})
// 	}

// 	return TeamCalendar
// }

func main() {

	// db.InitDatabase(DB_FILE)
	// db.ImportEmployeesFromCSV(DB_FILE, EMPLOYEES_CSV)
	// db.ImportTypesFromCSV(DB_FILE, EVENTTYPES_CSV, "eventTypes")
	// db.ImportTypesFromCSV(DB_FILE, DUTYTYPES_CSV, "dutyTypes")
	// db.ImportTablesFromCSV(DB_FILE, EVENTTABLE_CSV, "eventTable")
	// db.ImportTablesFromCSV(DB_FILE, DUTYTABLE_CSV, "dutyTable")

	// fmt.Println(createCalendar(5, 2024))

	sqlStmt := "SELECT * FROM eventTable"
	results, err := db.SelectFromTable(DB_FILE, sqlStmt)
	if err != nil {
		log.Fatalf("Query failed: %s", err)
	}

	for _, row := range results {
		fmt.Println(row)
	}
}
