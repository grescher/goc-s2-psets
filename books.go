package main

import (
	"math"

	"golang.org/x/exp/slices"
)

type AvgAgePerBook struct {
	BookTitle string
	AvgAge    int
}

func avgAgeOfReadersPerBook(users []User) (apb AvgAgePerBookSlice) {
	tmp := make(map[string]float64)
	for _, u := range users {
		for _, book := range u.Books {
			if _, ok := tmp[book]; !ok {
				tmp[book] = float64(u.Age)
			}
			tmp[book] = (float64(Age(u.Age)) + tmp[book]) / 2.0
		}
	}

	for k, v := range tmp {
		var ele AvgAgePerBook
		ele.BookTitle = k
		ele.AvgAge = int(math.Round(v))
		apb = append(apb, ele)
	}

	return apb
}

type AvgAgePerBookSlice []AvgAgePerBook

func (a AvgAgePerBookSlice) SortByAge() {
	b := []AvgAgePerBook(a)
	slices.SortFunc[AvgAgePerBook](b, func(x, y AvgAgePerBook) bool {
		return x.AvgAge > y.AvgAge
	})
}

func (a AvgAgePerBookSlice) NewTable(headers []string) (res Table) {
	res.Header = headers
	for _, ele := range a {
		field := make(RowField)
		field[res.Header[0]] = Name(ele.BookTitle).String()
		field[res.Header[1]] = Age(ele.AvgAge).String()

		res.Rows = append(res.Rows, field)
	}
	return res
}
