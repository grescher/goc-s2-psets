package main

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func main() {
	users := UserSlice(Users(data))

	// Practice #1: print the table with a users data.
	userHeaders := []string{"Name", "Age", "Active", "Mass", "Books"}
	PrintData(users, userHeaders)

	// Practice #2: print the average age of readers for each book.
	fmt.Println("\nBooks Stats:")
	avgAgeBooks := avgAgeOfReadersPerBook(users)
	averagePerBookHeaders := []string{"Book Title", "Avg.Age"}
	PrintData(avgAgeBooks, averagePerBookHeaders)

	// Practice #3:
	fmt.Println("\nSort users by the sum of the average age of the reader for each book they read:")
	sortUsersBySumOfAvgAge(users, avgAgeBooks)
	PrintData(users, userHeaders)

	fmt.Println("\nUser who have the mass as close to 80kg as possible:")
	if find, ok := users.FindMass(80); ok {
		PrintData(UserSlice{find}, userHeaders)
	}
}

func PrintData(data TablePrinter, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}

// Sorts the users by the sum of the average age for each book they read.
// Used in the 3rd practice.
func sortUsersBySumOfAvgAge(users []User, books []AvgAgePerBook) {
	ages := make(map[string]int)
	for _, book := range books {
		ages[book.BookTitle] = book.AvgAge
	}

	slices.SortFunc[User](users, func(x, y User) bool {
		var sumX, sumY int
		for _, book := range x.Books {
			sumX += ages[book]
		}
		for _, book := range y.Books {
			sumY += ages[book]
		}
		return sumX < sumY
	})
}
