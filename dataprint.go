package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	widthLimit    = 60
	separatorCol  = "|"
	separatorLine = "-"
)

type RowField map[string]string

type UsersToPrint struct {
	Header      []string
	Rows        []RowField
	ColumnWidth map[string]int
}

func PrintData(users []User, headers []string) {
	usersToPrint := NewUsersToPrint(users, headers)
	usersToPrint.TablePrint()
}

func NewUsersToPrint(users []User, headers []string) (res UsersToPrint) {
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
	res.setColumnWidth()

	return res
}

func (utp *UsersToPrint) setColumnWidth() {
	cw := make(map[string]int)
	for _, h := range utp.Header {
		cw[h] = utf8.RuneCountInString(h)
	}
	for _, fields := range utp.Rows {
		for _, h := range utp.Header {
			newlines := strings.Split(fields[h], "\n")
			for _, line := range newlines {
				width := utf8.RuneCountInString(line)
				if width > widthLimit {
					width = widthLimit
				}
				if width > cw[h] {
					cw[h] = width
				}
			}
		}
	}
	utp.ColumnWidth = cw
}

func (u *UsersToPrint) TablePrint() {
	u.printHeaders()
	u.printSeparator()
	u.printRows()
}

func (u *UsersToPrint) printHeaders() {
	var fields []string
	for _, h := range u.Header {
		lenOfH := utf8.RuneCountInString(h)
		left := (u.ColumnWidth[h]-lenOfH)/2 + lenOfH
		right := u.ColumnWidth[h] - left
		field := fmt.Sprintf(" %*s%*s ", left, h, right, "")
		fields = append(fields, field)
	}
	headerLine := strings.Join(fields, separatorCol)
	fmt.Printf("%[1]s%[2]s%[1]s\n", separatorCol, headerLine)
}

func (u *UsersToPrint) printSeparator() {
	var fields []string
	for _, h := range u.Header {
		line := strings.Repeat(separatorLine, u.ColumnWidth[h]+2)
		fields = append(fields, line)
	}
	line := strings.Join(fields, separatorCol)
	fmt.Printf("%[1]s%[2]s%[1]s\n", separatorCol, line)
}

func (u *UsersToPrint) printRows() {
	for _, row := range u.Rows {
		lines := getLinesToPrint(u.Header, u.ColumnWidth, row)

		for _, line := range lines {
			fmt.Printf("%[1]s%[2]s%[1]s\n", separatorCol, line)
		}
	}
}

func getLinesToPrint(header []string, columnWidth map[string]int, row map[string]string) (lines []string) {
	var isSingleLine bool
	for !isSingleLine {
		isSingleLine = true
		var chunks []string

		for _, h := range header {
			var ok bool
			var ch string

			ch, row[h], ok = chunk(row[h], columnWidth[h])
			if !ok {
				isSingleLine = false
			}
			chunks = append(chunks, ch)
		}
		lines = append(lines, strings.Join(chunks, separatorCol))
	}
	return lines
}

func chunk(str string, widthLimit int) (ch, newStr string, ok bool) {
	var b strings.Builder
	ok = true
	for i, width := 0, 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		width++
		if width == widthLimit {
			b.WriteRune(r)
			newStr, ok = str[i+size:], false
			ch = fmt.Sprintf(" %s ", b.String())
			return ch, newStr, ok
		}
		if r == '\n' {
			newStr, ok = str[i+size:], false
			break
		}
		b.WriteRune(r)
		i += size
	}
	ch = fmt.Sprintf(" %-*s ", widthLimit, b.String())
	return ch, newStr, ok
}
