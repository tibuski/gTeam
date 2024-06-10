package main

import (
	"github.com/tibuski/gTeam/db"
)

type Employee struct {
	Email   string
	Name    string
	Surname string
	Team    string
}

type Calendar struct {
	Email string
	Days  []string
}

var Employees = []Employee{}
var TeamCalendar = []Calendar{}

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

// func addEmployee(n, e, t string) Employee {
// 	human := Employee{Name: n, Email: e, Team: t}
// 	return human
// }

// func createCalendar(month, year int) []Calendar {
// 	d := daysOfTheMonth(month, year)
// 	for _, n := range Employees {
// 		TeamCalendar = append(TeamCalendar, Calendar{Name: n.Name, Days: d})
// 	}

// 	return TeamCalendar
// }

func main() {

	db.InitDatabase()
	db.ImportEmployeesFromCSV(db.CSV_FILE)

	// fmt.Println(createCalendar(5, 2024))
}
