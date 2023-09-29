package user

import (
	"fmt"
	"math"
	"practice/internal/table"
	"strings"

	"golang.org/x/exp/slices"
)

const massPrecision = 3

type User struct {
	Name        string   // uint8(length) + [length]byte
	Age         uint8    // 1 bit bool (active field) + 63 bit uint (age field)
	ActiveIndex uint8    // (see above)
	Mass        float64  // regular float64
	Books       []string // uint8(all books length) + [length]byte, all books come as a single comma-separated string
}

type Name string

func (n Name) String() string {
	return strings.Trim(fmt.Sprintf("%q", string(n)), "\"")
}

type Age uint8

func (a Age) String() string {
	return fmt.Sprintf("%d", uint8(a))
}

type ActiveIndex uint8

func (a ActiveIndex) String() string {
	if a > 0 {
		return "yes"
	}
	return "-"
}

type Mass float64

func (m Mass) String() string {
	res := fmt.Sprintf("%.*f", massPrecision, float64(m))
	res = strings.TrimRight(res, "0")
	if strings.HasSuffix(res, ".") {
		res = fmt.Sprint(res, "0")
	}
	res = fmt.Sprint(res, " kg")
	return res
}

type Books []string

func (b Books) String() string {
	var tmp []string
	for _, book := range b {
		tmp = append(tmp, fmt.Sprintf("%q", book))
	}
	return strings.Join(tmp, "\n")
}

type Slice []User

func (users Slice) NewTable(headers []string) (res table.Table) {
	res.Header = headers
	for _, user := range users {
		field := make(table.RowField)
		field[res.Header[0]] = Name(user.Name).String()
		field[res.Header[1]] = Age(user.Age).String()
		field[res.Header[2]] = ActiveIndex(user.ActiveIndex).String()
		field[res.Header[3]] = Mass(user.Mass).String()
		field[res.Header[4]] = Books(user.Books).String()

		res.Rows = append(res.Rows, field)
	}
	return res
}

func (u Slice) FindMass(m float64) (find User, ok bool) {
	users := make([]User, len(u))
	copy(users, u)
	slices.SortFunc[User](users, func(a, b User) bool {
		return a.Mass < b.Mass
	})

	idx, ok := slices.BinarySearchFunc[User, float64](users, m, func(u User, f float64) int {
		return int(math.Round(u.Mass - f))
	})
	if ok {
		find = users[idx]
	}
	return find, ok
}

func (u Slice) NumOfActiveUsers() (n int) {
	for _, user := range u {
		if user.ActiveIndex > 0 {
			n++
		}
	}
	return n
}

// Sorts the users by the sum of the average age for each book they read.
// Used in the 3rd practice.
func SortUsersBySumOfAvgAge(users []User, books []AvgAgePerBook) {
	ages := make(map[string]int)
	for _, book := range books {
		ages[book.BookTitle] = book.AvgAge
	}

	slices.SortFunc[User](users, func(x, y User) bool {
		var sumX, sumY int
		for _, book := range x.Books {
			sumX += ages[book]
		}
		for _, book := range y.Books {
			sumY += ages[book]
		}
		return sumX < sumY
	})
}
