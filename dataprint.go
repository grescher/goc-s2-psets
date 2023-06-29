package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	minNameWidth  = 6
	minAgeWidth   = 5
	activeWidth   = 8
	minMassWidth  = 6
	massPrecision = 8
)

func DataPrint(users []User) {
	ttyWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	tablePrint(users, ttyWidth)
}

func tablePrint(users []User, maxTableWidth int) {
	// Determine the width of the widest values.
	var (
		maxNameWidth int
		maxAgeWidth  int
		maxMassWidth int
	)

	for _, u := range users {
		if width := lenString(u.Name); width > maxNameWidth {
			maxNameWidth = width
		}
		if width := lenInt(u.Age); width > maxAgeWidth {
			maxAgeWidth = width
		}
		if width := lenFloat64(u.Mass, massPrecision); width > maxMassWidth {
			maxMassWidth = width
		}
	}

	fmt.Printf("| %v | %v | %v | %v |\n", "Name", "Age", "Active", "Mass")
	for i := 0; i < maxTableWidth; i++ {
		fmt.Print("-")
	}
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

func lenFloat64(num float64, precision int) (width int) {
	s := fmt.Sprintf("%.*f", precision, num)
	s = strings.TrimRight(s, "0")
	width = len(s)
	if r, _ := utf8.DecodeLastRuneInString(s); r == '.' {
		return width + 1
	}
	return width
}
