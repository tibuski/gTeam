package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tibuski/gTeam/db"
)

type calendarDay struct {
	code    string `default:"0"`
	bgColor string `default:"000000"`
}

type Calendar struct {
	name           string
	surname        string
	days           []calendarDay
	employeeNumber int
	email          string
	team           string
}

const DB_FILE = "./db/gTeam.db"
const EMPLOYEES_CSV = "./db/people.csv"
const EVENTTYPES_CSV = "./db/eventTypes.csv"
const EVENTTABLE_CSV = "./db/eventTable.csv"
const DUTYTYPES_CSV = "./db/dutyTypes.csv"
const DUTYTABLE_CSV = "./db/dutyTable.csv"

func daysOfTheMonth(month int, year int) []calendarDay {
	var days []calendarDay
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
		if t.Weekday() == 6 || t.Weekday() == 0 {
			days = append(days, calendarDay{code: "W", bgColor: "000000"})
		} else {
			days = append(days, calendarDay{code: "0", bgColor: "000000"})
		}
	}
	return days
}

func createCalendar(database *sql.DB, month int, year int) []Calendar {
	// Insert days of the month in TeamCalendar
	d := daysOfTheMonth(month, year)
	p, _ := db.SelectFromPeople(database, "%")
	e, _ := db.SelectPeopleAllEvents(database, month, year)

	var TeamCalendar = []Calendar{}
	for n := range len(p) {
		TeamCalendar = append(TeamCalendar, Calendar{name: p[n].Name, surname: p[n].Surname, days: d, email: p[n].Email, employeeNumber: p[n].EmployeeNumber, team: p[n].Team})
	}
	// Insert all events in TeamCalendar

	fmt.Println(d)
	fmt.Println(p)
	fmt.Println(e)

	return TeamCalendar
}

func main() {
	database := db.OpenDatabase(DB_FILE)
	// db.InitDatabase(database)
	// db.ImportEmployeesFromCSV(database, EMPLOYEES_CSV)
	// db.ImportTypesFromCSV(database, EVENTTYPES_CSV, "eventTypes")
	// db.ImportTypesFromCSV(database, DUTYTYPES_CSV, "dutyTypes")
	// db.ImportTablesFromCSV(database, EVENTTABLE_CSV, "eventTable")
	// db.ImportTablesFromCSV(database, DUTYTABLE_CSV, "dutyTable")
	fmt.Println(createCalendar(database, 5, 2024))

}
