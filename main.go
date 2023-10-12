package main

import (
	"log"
	"practice/internal/storage"
	"practice/internal/tui"
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

	tui.Prompt(strg, &users)
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
