package main

import (
	"fmt"
	"log"
	"practice/internal/storage"
	"practice/internal/table"
	"practice/internal/user"
)

func main() {
	storage, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	users, err := user.Decode(storage.Reader())
	if err != nil {
		log.Fatal(err)
	}
	defer saveSnapshot(storage, &users)

	userHeaders := []string{"Name", "Age", "Active", "Mass", "Books"}

	user.SortUsersBySumOfAvgAge(users, user.AvgAgeOfReadersPerBook(users))
	PrintData(user.Slice(users), userHeaders)
	fmt.Println("Number of active users:", user.Slice(users).NumOfActiveUsers())
}

func saveSnapshot(storage *storage.Storage, users *[]user.User) {
	data := user.Encode(*users)
	storage.Write(data)
}

func PrintData(data table.Printer, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}
