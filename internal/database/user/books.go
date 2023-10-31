package user

import (
	"math"
	"practice/internal/database/tui/table"

	"golang.org/x/exp/slices"
)

type AvgAgePerBook struct {
	BookTitle string
	AvgAge    int
}

func AvgAgeOfReadersPerBook(users []User) (apb AvgAgePerBookSlice) {
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

// NewTable method satisfies the table.Printer interface.
// It converts the slice data to strings and as a result creates a new table.Table object.
func (a AvgAgePerBookSlice) NewTable(headers []string) (res table.Table) {
	res.Headers = headers
	for _, ele := range a {
		// Create a new row and fill it with values for each column.
		row := make(table.Row)
		row[res.Headers[0]] = Name(ele.BookTitle).String()
		row[res.Headers[1]] = Age(ele.AvgAge).String()

		res.Rows = append(res.Rows, row)
	}
	return res
}
