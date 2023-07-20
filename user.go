package main

import (
	"fmt"
	"strings"
)

const (
	massPrecision = 3
	ozPerKg       = 35.273962
)

type User struct {
	Name   string
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
		field[res.Header[1]] = Age(user.Age).String()
		field[res.Header[2]] = Active(user.Active).String()
		field[res.Header[3]] = Mass(user.Mass).String()
		field[res.Header[4]] = Books(user.Books).String()

		res.Rows = append(res.Rows, field)
	}
	return res
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
