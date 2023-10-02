package main

import (
	"fmt"
	"log"
	"practice/internal/storage"
	"practice/internal/table"
	"practice/internal/user"
)

func main() {
	strg, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer closeStorage(strg)

	users, err := user.Decode(strg.Reader())
	if err != nil {
		log.Fatal(err)
	}
	defer saveSnapshot(strg, &users)

	userHeaders := []string{"Name", "Age", "Active", "Mass", "Books"}

	user.SortUsersBySumOfAvgAge(users, user.AvgAgeOfReadersPerBook(users))
	PrintData(user.Slice(users), userHeaders)
	fmt.Println("Number of active users:", user.Slice(users).NumOfActiveUsers())
}

func closeStorage(strg *storage.Storage) {
	if err := strg.Close(); err != nil {
		log.Fatal("closeStorage: ", err)
	}
	log.Println("Done. Bye.")
}

func saveSnapshot(strg *storage.Storage, users *[]user.User) {
	log.Print("Saving snapshot... ")
	if err := strg.SaveSnapshot(users); err != nil {
		log.Fatal("saveSnapshot: ", err)
	}
}

func PrintData(data table.Printer, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}
