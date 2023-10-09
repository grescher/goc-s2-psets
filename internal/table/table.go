package table

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	widthLimit    = 60
	separatorCol  = "|"
	separatorRow  = "+"
	separatorLine = "-"
)

type RowField map[string]string

type Table struct {
	Header      []string
	Rows        []RowField
	ColumnWidth map[string]int
}

type Printer interface {
	NewTable([]string) Table
}

func PrintData(data Printer, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}

func (t *Table) Print() {
	t.setColumnWidth()
	t.printHeaders()
	t.printSeparator(separatorCol)
	t.printRows()
}

func (t *Table) setColumnWidth() {
	cw := make(map[string]int)
	for _, h := range t.Header {
		cw[h] = utf8.RuneCountInString(h)
	}
	for _, fields := range t.Rows {
		for _, h := range t.Header {
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
	t.ColumnWidth = cw
}

func (t *Table) printHeaders() {
	var fields []string
	for _, h := range t.Header {
		lenOfH := utf8.RuneCountInString(h)
		left := (t.ColumnWidth[h]-lenOfH)/2 + lenOfH
		right := t.ColumnWidth[h] - left
		field := fmt.Sprintf(" %*s%*s ", left, h, right, "")
		fields = append(fields, field)
	}
	headerLine := strings.Join(fields, separatorCol)
	fmt.Printf("%[1]s%[2]s%[1]s\n", separatorCol, headerLine)
}

func (t *Table) printSeparator(sep string) {
	var fields []string
	for _, h := range t.Header {
		line := strings.Repeat(separatorLine, t.ColumnWidth[h]+2)
		fields = append(fields, line)
	}
	line := strings.Join(fields, sep)
	fmt.Printf("%[1]s%[2]s%[1]s\n", separatorCol, line)
}

func (t *Table) printRows() {
	for _, row := range t.Rows {
		lines := getLinesToPrint(t.Header, t.ColumnWidth, row)

		for _, line := range lines {
			fmt.Printf("%[1]s%[2]s%[1]s\n", separatorCol, line)
		}
		t.printSeparator(separatorRow)
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
