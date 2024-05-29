package main

import (
	"fmt"
	"time"

	_ "github.com/tibuski/gTeam/db"
)

type Employee struct {
	Name  string
	Email string
	Team  string
}

type Calendar struct {
	Name string
	Days []string
}

var Employees = []Employee{}
var TeamCalendar = []Calendar{}

func daysOfTheMonth(month int, year int) []string {
	var days []string
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
		if t.Weekday() == 6 || t.Weekday() == 0 {
			days = append(days, "W")
		} else {
			days = append(days, "0")
		}
	}
	return days
}

func addEmployee(n, e, t string) Employee {
	human := Employee{Name: n, Email: e, Team: t}
	return human
}

func createCalendar(month, year int) []Calendar {
	d := daysOfTheMonth(month, year)
	for _, n := range Employees {
		TeamCalendar = append(TeamCalendar, Calendar{Name: n.Name, Days: d})
	}

	return TeamCalendar
}

func main() {
	db.initDatabase()

	Employees = append(Employees, addEmployee("Laurent Brichet", "laurent@brichet.be", "CVVS"))
	Employees = append(Employees, addEmployee("Gerard Pascal", "gerard@proximus.lu", "BCNS"))
	Employees = append(Employees, addEmployee("Pol Dupont", "pol.dupont@proximus.lu", "CVVS"))
	Employees = append(Employees, addEmployee("Roger Jaco", "Roro.Gege@proximus.lu", "BCNS"))

	// fmt.Println(Employees)
	// fmt.Println(len(NumberOfDays))
	// fmt.Println(NumberOfDays)
	fmt.Println(createCalendar(5, 2024))
}
