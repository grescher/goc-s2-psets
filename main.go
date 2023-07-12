package main

import (
	"fmt"
	"math"
)

func main() {
	users := UserSlice(Users())
	PrintData(users, UserHeaders)
	fmt.Println("\nBooks Stats:")

	books := avgAgeOfReadersPerBook(users)
	PrintData(books, APBHeaders)
}

func PrintData(data TablePrinter, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}

type AgePerBook map[string]int

var APBHeaders = []string{"Name", "Avg. Age"}

func avgAgeOfReadersPerBook(users UserSlice) AgePerBook {
	tmp := make(map[string]float64)
	for _, u := range users {
		for _, book := range u.Books {
			if _, ok := tmp[book]; !ok {
				tmp[book] = float64(u.Age)
			}
			tmp[book] = (float64(Age(u.Age).Verify()) + tmp[book]) / 2.0
		}
	}

	apb := make(map[string]int)
	for k, v := range tmp {
		apb[k] = int(math.Round(v))
	}

	return apb
}

func (a AgePerBook) NewTable(headers []string) (res Table) {
	res.Header = headers
	for name, age := range a {
		field := make(RowField)
		field[res.Header[0]] = Name(name).String()
		field[res.Header[1]] = Age(age).String()

		res.Rows = append(res.Rows, field)
	}
	return res
}
