package main

import (
	"fmt"
	"log"

	"golang.org/x/exp/slices"
)

func main() {
	users, err := Users(Reader())
	if err != nil {
		log.Fatal(err)
	}

	userHeaders := []string{"Name", "Age", "Active", "Mass", "Books"}

	sortUsersBySumOfAvgAge(users, avgAgeOfReadersPerBook(users))
	PrintData(UserSlice(users), userHeaders)
	fmt.Println("Number of active users:", UserSlice(users).NumOfActiveUsers())
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
