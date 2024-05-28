package main

import "fmt"

type Employee struct {
	Name  string
	Email string
	Team  string
}

var Employees = []Employee{}

func addEmployee(n, e, t string) Employee {
	human := Employee{Name: n, Email: e, Team: t}
	return (human)
}

func main() {
	Employees = append(Employees, addEmployee("Laurent", "l@b.be", "CVVS"))
	Employees = append(Employees, addEmployee("Gerard", "t@po.fr", "BCNS"))

	fmt.Println(Employees)
}
