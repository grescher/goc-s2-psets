package main

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	massPrecision = 3
	ozPerKg       = 35.273962
)

type TypeOfUser int8

const (
	Undefined TypeOfUser = iota
	Organizer
	Critic
	NewbieReader
	CasualReader
	NonFictionBuff
	PoetryLover
	SciFiFan
)

func (t TypeOfUser) String() (s string) {
	switch t {
	case Undefined:
		s = "Undefined"
	case Organizer:
		s = "Organizer"
	case Critic:
		s = "Critic"
	case NewbieReader:
		s = "NewbieReader"
	case CasualReader:
		s = "CasualReader"
	case NonFictionBuff:
		s = "NonFictionBuff"
	case PoetryLover:
		s = "PoetryLover"
	case SciFiFan:
		s = "SciFiFan"
	}
	return s
}

type User struct {
	Name   string
	Type   TypeOfUser
	Age    int
	Active bool
	Mass   float64
	Books  []string
}

func Users(data []User) []User {
	return data
}

type UserSlice []User

func (users UserSlice) NewTable(headers []string) (res Table) {
	res.Header = headers
	for _, user := range users {
		field := make(RowField)
		field[res.Header[0]] = Name(user.Name).String()
		field[res.Header[1]] = user.Type.String()
		field[res.Header[2]] = Age(user.Age).String()
		field[res.Header[3]] = Active(user.Active).String()
		field[res.Header[4]] = Mass(user.Mass).String()
		field[res.Header[5]] = Books(user.Books).String()

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

type Name string

func (n Name) String() string {
	return strings.Trim(fmt.Sprintf("%q", string(n)), "\"")
}

type Age int

func (a Age) Verify() Age {
	if a < 0 {
		return 100 + a
	}
	return a
}

func (a Age) String() string {
	return fmt.Sprintf("%d", int(a.Verify()))
}

type Active bool

func (a Active) String() string {
	if a {
		return "yes"
	}
	return "-"
}

type Mass float64

func (m Mass) String() string {
	var abbr string
	switch {
	case m > 0.0009 && m < 1:
		abbr = "qq" // quintal/centner (100kg)
	case m > 620:
		abbr = "oz" // ounces
	default:
		abbr = "kg"
	}
	res := fmt.Sprintf("%.*f", massPrecision, float64(m))
	res = strings.TrimRight(res, "0")
	if strings.HasSuffix(res, ".") {
		res = fmt.Sprint(res, "0")
	}
	res = fmt.Sprint(res, " ", abbr)
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
