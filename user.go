package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

const massPrecision = 3

const (
	ActiveMask    = 1 << 63
	AgeMask       = math.MaxUint64 ^ ActiveMask
	maxNumOfUsers = 8
	kgPerOz       = 0.0283495
	kgPerQq       = 100.0
)

type User struct {
	Name        string   // uint8(length) + [length]byte
	Age         uint8    // 1 bit bool (active field) + 63 bit uint (age field)
	ActiveIndex uint8    // (see above)
	Mass        float64  // regular float64
	Books       []string // uint8(all books length) + [length]byte, all books come as a single comma-separated string
}

func UsersDecode(r io.Reader) (out []User, err error) {
	for err != io.EOF && len(out) < maxNumOfUsers {
		var nameLength uint8
		if err = binary.Read(r, binary.BigEndian, &nameLength); err != nil {
			break
		}
		name := make([]byte, nameLength)
		if err = binary.Read(r, binary.BigEndian, &name); err != nil {
			break
		}
		var activeAndAge uint64
		if err = binary.Read(r, binary.BigEndian, &activeAndAge); err != nil {
			break
		}
		var active uint8
		if activeAndAge&ActiveMask > 0 {
			active = 1
		}
		var mass float64
		if err = binary.Read(r, binary.BigEndian, &mass); err != nil {
			break
		}
		var booksLen uint8
		if err = binary.Read(r, binary.BigEndian, &booksLen); err != nil {
			break
		}
		books := make([]byte, booksLen)
		if err = binary.Read(r, binary.BigEndian, &books); err != nil {
			break
		}

		var user User
		user.Name = string(name)
		user.ActiveIndex = active << len(out)
		user.Age = uint8(activeAndAge & AgeMask)
		user.Mass = verifyMass(mass)
		user.Books = strings.Split(string(books), ",")
		out = append(out, user)
	}
	if len(out) != maxNumOfUsers && err != io.EOF {
		return nil, err
	}
	return out, nil
}

func verifyMass(m float64) float64 {
	switch {
	case m > 0.0009 && m < 1: // quintals to kg
		m = m * kgPerQq
	case m > 620:
		m = m * kgPerOz // ounces to kg
	}
	return m
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

type UserSlice []User

func (users UserSlice) NewTable(headers []string) (res Table) {
	res.Header = headers
	for _, user := range users {
		field := make(RowField)
		field[res.Header[0]] = Name(user.Name).String()
		field[res.Header[1]] = Age(user.Age).String()
		field[res.Header[2]] = ActiveIndex(user.ActiveIndex).String()
		field[res.Header[3]] = Mass(user.Mass).String()
		field[res.Header[4]] = Books(user.Books).String()

		res.Rows = append(res.Rows, field)
	}
	return res
}

func (u UserSlice) FindMass(m float64) (find User, ok bool) {
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

func (u UserSlice) NumOfActiveUsers() (n int) {
	for _, user := range u {
		if user.ActiveIndex > 0 {
			n++
		}
	}
	return n
}

// Sorts the users by the sum of the average age for each book they read.
// Used in the 3rd practice.
func sortUsersBySumOfAvgAge(users []User, books []AvgAgePerBook) {
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
