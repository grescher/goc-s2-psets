package main

import "fmt"

func main() {
	users := UserSlice(Users())
	PrintData(users, UserHeaders)
	fmt.Println("\nBooks Stats:")

	books := avgAgeOfReadersPerBook(users)
	books.SortByAge()
	PrintData(books, APBHeaders)
}

func PrintData(data TablePrinter, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}
