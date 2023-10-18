package main

import (
	"log"
	"practice/internal/storage"
	"practice/internal/tcp"
	"practice/internal/user"
)

func main() {
	// Open/create a storage.
	strg, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer closeStorage(strg)

	// Read the data from the storage.
	users, err := user.Decode(strg.Reader())
	if err != nil {
		log.Fatal(err)
	}
	defer saveSnapshot(strg, &users)

	c := make(chan int)
	// Start a TCP server.
	go tcp.Server(c, strg, &users)
	// Start a TCP client.
	go tcp.Client(c)

	<-c
	<-c
	// Show the text user interface prompt.
	// tui.Prompt(os.Stdin, os.Stdout, strg, &users)
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
