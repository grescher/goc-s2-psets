package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	minNameWidth  = 4
	minAgeWidth   = 3
	activeWidth   = 6
	minMassWidth  = 4
	massPrecision = 8
	maxFieldWidth = 60
)

func DataPrint(users []User) {
	// ValueLen contains length of stringified value for each field of the users slice.
	type ValueLen struct {
		Name     int
		Age      int
		Mass     int
		MassPrec int
	}

	// Get all the width values for each item of the users slice.
	var valLen []ValueLen
	for _, u := range users {
		var vl ValueLen
		vl.Name = lenString(u.Name)
		vl.Age = lenInt(u.Age)
		vl.Mass, vl.MassPrec = lenFloat64(u.Mass, massPrecision)

		valLen = append(valLen, vl)
	}

	// Determine and set the max width for each field.
	var fieldWidth struct {
		Name   int
		Age    int
		Active int
		Mass   int
	}

	fieldWidth.Name = minNameWidth
	fieldWidth.Age = minAgeWidth
	fieldWidth.Active = activeWidth
	fieldWidth.Mass = minMassWidth

	for _, vl := range valLen {
		if fieldWidth.Name < vl.Name {
			if vl.Name > maxFieldWidth {
				fieldWidth.Name = maxFieldWidth
			} else {
				fieldWidth.Name = vl.Name
			}
		}
		if fieldWidth.Age < vl.Age {
			if vl.Age > maxFieldWidth {
				fieldWidth.Age = maxFieldWidth
			} else {
				fieldWidth.Age = vl.Age
			}
		}
		if fieldWidth.Mass < vl.Mass {
			if vl.Mass > maxFieldWidth {
				fieldWidth.Mass = maxFieldWidth
			} else {
				fieldWidth.Mass = vl.Mass
			}
		}
	}

	// Print header of the table.
	fmt.Printf(
		"| %*s%*s | %*s%*s | %s | %*s%*s |\n",

		(fieldWidth.Name-minNameWidth)/2, "Name",
		fieldWidth.Name-(fieldWidth.Name-minNameWidth)/2-minNameWidth, " ",

		(fieldWidth.Age-minAgeWidth)/2, "Age",
		fieldWidth.Age-(fieldWidth.Age-minAgeWidth)/2-minAgeWidth, " ",

		"Active",

		(fieldWidth.Mass-minMassWidth)/2, "Mass",
		fieldWidth.Mass-(fieldWidth.Mass-minMassWidth)/2-minMassWidth, " ",
	)

	for _, user := range users {
		fmt.Printf("%[2]*[1]v (%v) |", user.Name, len(user.Name)+1, len(user.Name))
		fmt.Printf("%[2]*[1]v |", user.Age, len(strconv.Itoa(user.Age))+1)
		var strActive string
		if user.Active {
			strActive = "true"
		} else {
			strActive = "false"
		}
		fmt.Printf("%[2]*[1]t |", user.Active, len(strActive)+1)
		fmt.Printf("%#g\n", user.Mass)
	}
}

func lenString(s string) (width int) {
	width, _ = fmt.Fprintf(io.Discard, "%q", s)
	return width - 2 // Extract the quotation marks
}

func lenInt(num int) (width int) {
	if num == 0 {
		return 1
	}
	if num < 0 {
		width++
	}
	for num != 0 {
		num /= 10
		width++
	}
	return width
}

func lenFloat64(num float64, maxPrecision int) (width, precision int) {
	s := fmt.Sprintf("%.*f", maxPrecision, num)
	s = strings.TrimRight(s, "0")
	width = len(s)

	i := strings.IndexRune(s, '.')
	if i == width-1 || i <= 0 { // The trailing decimal zero remains
		width++
		precision = 1
	} else {
		precision = (width - 1) - i
	}

	return width, precision
}
