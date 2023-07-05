package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	minNameWidth  = 4 // width of string "Name"
	minAgeWidth   = 3 //             ... "Age"
	activeWidth   = 6 //             ... "Active"
	minMassWidth  = 4 //             ... "Mass"
	massPrecision = 8
	widthLimit    = 60
)

// PrintFloat contains width and precision to print
type PrintFloat struct {
	Width     int
	Precision int
}

// FieldWidth contains widths for string values of a User type.
type FieldWidth struct {
	Name   int
	Age    int
	Active int
	Mass   PrintFloat
}

type LineToPrint struct {
	Name   string
	Age    string
	Active string
	Mass   string
}

func DataPrint(users []User) {
	// Determine and set the max width to be printed for each field.
	widthOfField := getMaxFieldWidths(users)

	// Print header of the table.
	left := (widthOfField.Name-minNameWidth)/2 + minNameWidth
	right := widthOfField.Name - left
	fmt.Printf("| %*s%*s ", left, "Name", right, "")

	left = (widthOfField.Age-minAgeWidth)/2 + minAgeWidth
	right = widthOfField.Age - left
	fmt.Printf("| %*s%*s ", left, "Age", right, "")

	fmt.Printf("| %s ", "Active")

	left = (widthOfField.Mass.Width-minMassWidth)/2 + minMassWidth
	right = widthOfField.Mass.Width - left
	fmt.Printf("| %*s%*s |\n", left, "Mass", right, "")

	// Print a header underline.
	fmt.Printf(
		"|%s|%s|%s|%s|\n",
		strings.Repeat("-", widthOfField.Name+2),
		strings.Repeat("-", widthOfField.Age+2),
		strings.Repeat("-", widthOfField.Active+2),
		strings.Repeat("-", widthOfField.Mass.Width+2),
	)

	// Print the table rows.
	for _, user := range users {
		lines := getLinesToPrint(user, widthOfField)
		printLines(lines)
	}
}

// getMaxFieldWidths returns the width of the widest columns of the table to be printed.
func getMaxFieldWidths(users []User) (mfw FieldWidth) {
	mfw.Name = minNameWidth
	mfw.Age = minAgeWidth
	mfw.Active = activeWidth
	mfw.Mass.Width = minMassWidth

	for _, u := range users {
		w := getFieldWidths(u)
		if mfw.Name < w.Name {
			mfw.Name = w.Name
		}
		if mfw.Age < w.Age {
			mfw.Age = w.Age
		}
		if mfw.Mass.Width < w.Mass.Width {
			mfw.Mass.Width = w.Mass.Width
		}
	}
	if mfw.Name > widthLimit {
		mfw.Name = widthLimit
	}
	if mfw.Age > widthLimit {
		mfw.Age = widthLimit
	}
	if mfw.Mass.Width > widthLimit {
		mfw.Mass.Width = widthLimit
	}
	return mfw
}

func getLinesToPrint(user User, fwidth FieldWidth) (lines []LineToPrint) {
	nameStr := fmt.Sprintf("%q", user.Name)
	nameStr = strings.Trim(nameStr, "\"")
	ageStr := fmt.Sprintf("%d", user.Age)
	massWidth, massPrec := lenFloat64(user.Mass, massPrecision)
	massStr := fmt.Sprintf("%*.*f", massWidth, massPrec, user.Mass)
	activeStr := "-"
	if user.Active {
		activeStr = "yes"
	}

	var isSingleLine bool
	for !isSingleLine {
		isSingleLine = true
		var line LineToPrint
		var ok bool

		if line.Name, ok = setLine(&nameStr, fwidth.Name); !ok {
			isSingleLine = false
		}
		if line.Age, ok = setLine(&ageStr, fwidth.Age); !ok {
			isSingleLine = false
		}
		if line.Active, ok = setLine(&activeStr, fwidth.Active); !ok {
			isSingleLine = false
		}
		if line.Mass, ok = setLine(&massStr, fwidth.Mass.Width); !ok {
			isSingleLine = false
		}
		lines = append(lines, line)
	}
	return lines
}

func setLine(s *string, widthOfField int) (line string, isSingleLine bool) {
	line = *s
	width := utf8.RuneCountInString(line)
	if width > widthOfField {
		*s = line[widthOfField:]
		return line[:widthOfField], false
	}
	line = fmt.Sprintf("%-*s", widthOfField, *s)
	*s = ""
	return line, true
}

// printLinew prints the User's record row by row (if width of one of the fields
// exceeds the maximum column width.)
func printLines(lines []LineToPrint) {
	for _, l := range lines {
		fmt.Printf(
			"| %s | %s | %s | %s |\n",
			l.Name, l.Age, l.Active, l.Mass,
		)
	}
}

// getFieldWidths returns the width of each field of the User structure values
// to be converted in a quoted string.
func getFieldWidths(u User) (fw FieldWidth) {
	fw.Name = lenString(u.Name)
	fw.Age = lenInt(u.Age)
	if u.Active {
		fw.Active = activeWidth
	}
	fw.Mass.Width, fw.Mass.Precision = lenFloat64(u.Mass, massPrecision)

	return fw
}

// lenString returns the width of a printed string.
func lenString(s string) (width int) {
	width, _ = fmt.Fprintf(io.Discard, "%q", s)
	return width - 2 // Extract the width of quotation marks
}

// lenInt returns the width of a printed integer.
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

// lenFloat64 returns the width of a printed float64 and its precision.
// If num has no digits after the decimal point, the precision will be 1 to keep
// the zero in the tenths place.
// (i.e. lenFloat64(7) will return (3, 1) to get the printed string "7.0").
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
