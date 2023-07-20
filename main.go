package main

import "fmt"

func main() {
	users := UserSlice(Users(data))

	// Practice #1: print the table with a users data.
	userHeaders := []string{"Name", "Age", "Active", "Mass", "Books"}
	PrintData(users, userHeaders)

	// Practice #2: print the average age of readers for each book.
	fmt.Println()
	fmt.Println("Books Stats:")
	books := avgAgeOfReadersPerBook(users)
	books.SortByAge()
	averagePerBookHeaders := []string{"Book Title", "Avg.Age"}
	PrintData(books, averagePerBookHeaders)

}

func PrintData(data TablePrinter, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}
