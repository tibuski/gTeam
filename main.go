package main

import (
	"fmt"
	"time"
)

type Employee struct {
	Name  string
	Email string
	Team  string
}

type Calendar struct {
	Name string
	Days []int
}

var Employees = []Employee{}

func daysOfTheMonth(month time.Month, year int) []int {
	var days []int
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(start)
	end := start.AddDate(0, 1, -1)
	fmt.Println(end)
	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
		days = append(days, int(t.Weekday()))
	}
	return days
}

func addEmployee(n, e, t string) Employee {
	human := Employee{Name: n, Email: e, Team: t}
	return (human)
}

func main() {

	Employees = append(Employees, addEmployee("Laurent", "laurent@brichet.be", "CVVS"))
	Employees = append(Employees, addEmployee("Gerard", "gerard@proximus.lu", "BCNS"))
	Employees = append(Employees, addEmployee("Pol", "pol.dupont@proximus.lu", "CVVS"))
	Employees = append(Employees, addEmployee("Roger", "Roro.Gege@proximus.lu", "BCNS"))

	NumberOfDays := daysOfTheMonth(5, 2024)
	// fmt.Println(Employees)
	fmt.Println(len(NumberOfDays))
	fmt.Println(NumberOfDays)

}
