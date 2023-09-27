package main

import (
	"fmt"
	"log"
	"practice/internal/storage"
)

func main() {
	storage, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	users, err := UsersDecode(storage.Reader())
	if err != nil {
		log.Fatal(err)
	}
	// defer storage.Write(users)

	userHeaders := []string{"Name", "Age", "Active", "Mass", "Books"}

	sortUsersBySumOfAvgAge(users, avgAgeOfReadersPerBook(users))
	PrintData(UserSlice(users), userHeaders)
	fmt.Println("Number of active users:", UserSlice(users).NumOfActiveUsers())
}

func PrintData(data TablePrinter, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}
