package table

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	// Limit for the column width.
	widthLimit = 38

	// Symbols to print a table border.
	// Regular rows:
	symLine      = "─"
	symLineLeft  = "├"
	symLineRight = "┤"
	// Header rows:
	symDoubleLine      = "═"
	symLineDoubleLeft  = "╞"
	symDoubleCross     = "╪"
	symLineDoubleRight = "╡"
	// Columns:
	symColumnLine = "│"
	symColumnUp   = "┬"
	symCross      = "┼"
	symColumnDown = "┴"
	// Corners:
	symCornerUpLeft    = "╭"
	symCornerUpRight   = "╮"
	symCornerDownLeft  = "╰"
	symCornerDownRight = "╯"
)

// Row contains values for each column in this row.
type Row map[string]string

// ColumnWidth determines width for each column.
func (r *Row) ColumnWidths(headers []string) map[string]int {
	cw := make(map[string]int)
	for _, colName := range headers {
		// cell may have more than one line
		lines := strings.Split((*r)[colName], "\n")

		for _, line := range lines {
			width := utf8.RuneCountInString(line)

			if width > widthLimit {
				// the width of the column cannot exceed the width limit.
				width = widthLimit
			}
			if width > cw[colName] {
				cw[colName] = width
			}
		}
	}
	return cw
}

// Table type is used to print data objects that satisfy the Printer interface.
type Table struct {
	// Headers contains the names of the columns that will be printed in the header.
	Headers []string

	// Rows contains data as slice of maps containing values for each cell of the row,
	// i.e. []Row[columnName]cellValue.
	Rows []Row

	// Column Width contains the maximum width value for each column of the table,
	// does not include leading and trailing spaces.
	ColumnWidth map[string]int
}

type Printer interface {
	NewTable([]string) Table
}

func PrintData(w io.Writer, data Printer, headers []string) {
	table := data.NewTable(headers)
	table.Print(w)
}

func (t *Table) Print(w io.Writer) {
	t.setColumnWidth()
	t.printLine(w, symCornerUpLeft, symLine, symColumnUp, symCornerUpRight)
	t.printHeaders(w)
	t.printLine(w, symLineDoubleLeft, symDoubleLine, symDoubleCross, symLineDoubleRight)
	t.printRows(w)
	t.printLine(w, symCornerDownLeft, symLine, symColumnDown, symCornerDownRight)
}

// setColumnWidth scans each row of the table and determines maximum width
// for each column.
func (t *Table) setColumnWidth() {
	columnWidth := make(map[string]int)
	for _, h := range t.Headers {
		columnWidth[h] = utf8.RuneCountInString(h)
	}
	for _, row := range t.Rows {
		rowWidths := row.ColumnWidths(t.Headers)
		for _, colName := range t.Headers {
			if rowWidths[colName] > columnWidth[colName] {
				columnWidth[colName] = rowWidths[colName]
			}
		}
	}
	t.ColumnWidth = columnWidth
}

// printHeaders prints a line of the table with the column names.
func (t *Table) printHeaders(w io.Writer) {
	var cells []string
	for _, hName := range t.Headers {
		hNameLen := utf8.RuneCountInString(hName)

		// Construct a string with the column name centered in the header cell.
		columnWidth := t.ColumnWidth[hName]
		leftSpace := (columnWidth-hNameLen)/2 + hNameLen
		rightSpace := columnWidth - leftSpace
		cellContent := fmt.Sprintf(" %*s%*s ", leftSpace, hName, rightSpace, "")
		cells = append(cells, cellContent)
	}
	// Join cells separated by a column separation symbol into one line.
	headerLine := strings.Join(cells, symColumnLine)
	// Print the line with outer border symbols.
	fmt.Fprintf(w, "%[1]s%[2]s%[1]s\n", symColumnLine, headerLine)
}

// printLine prints a line separator.
func (t *Table) printLine(w io.Writer, left, dash, column, right string) {
	var cells []string
	for _, hName := range t.Headers {
		line := strings.Repeat(dash, t.ColumnWidth[hName]+2)
		cells = append(cells, line)
	}
	// Join line with a symbol of column separation.
	line := strings.Join(cells, column)
	// Print line with edge symbols.
	fmt.Fprintf(w, "%s%s%s\n", left, line, right)
}

// printRows prints rows of the table separated by a line.
func (t *Table) printRows(w io.Writer) {
	for i, row := range t.Rows {
		// One table row may consist of several lines.
		lines := getLinesToPrint(t.Headers, t.ColumnWidth, row)
		// Print lines of the current row.
		for _, line := range lines {
			fmt.Fprintf(w, "%[1]s%[2]s%[1]s\n", symColumnLine, line)
		}
		// Print a line separator between rows.
		if i != len(t.Rows)-1 {
			t.printLine(w, symLineLeft, symLine, symCross, symLineRight)
		}
	}
}

// getLinesToPrint returns a slice of lines that make up a table row.
func getLinesToPrint(headers []string, columnWidth map[string]int, row map[string]string) (lines []string) {
	var isSingleLine bool
	for !isSingleLine {
		isSingleLine = true
		var cellsOfLine []string

		for _, h := range headers {
			cell, remainderOfContent, isRemainderLeft := getLineOfCell(row[h], columnWidth[h])
			row[h] = remainderOfContent
			if isRemainderLeft {
				isSingleLine = false
			}
			cellsOfLine = append(cellsOfLine, cell)
		}
		line := strings.Join(cellsOfLine, symColumnLine)
		lines = append(lines, line)
	}
	return lines
}

// getLineOfCell returns the line of cell that is prepared for printing, the remaining string
// of the cell content if any, and a flag if there is any remaining content.
func getLineOfCell(str string, widthLimit int) (loc, strRemainder string, hasRemainder bool) {
	var b strings.Builder
	for i, width := 0, 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		width++
		if width == widthLimit {
			b.WriteRune(r)
			strRemainder = str[i+size:]
			hasRemainder = (len(strRemainder) > 0)
			break
		}
		if r == '\n' {
			strRemainder = str[i+size:]
			hasRemainder = (len(strRemainder) > 0)
			break
		}
		b.WriteRune(r)
		i += size
	}
	// Assemble the line of cell.
	loc = fmt.Sprintf(" %-*s ", widthLimit, b.String())
	return loc, strRemainder, hasRemainder
}
