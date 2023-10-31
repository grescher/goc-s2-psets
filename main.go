package main

import (
	"log"
	"practice/internal/database/storage"
	"practice/internal/database/tcp"
	"practice/internal/database/user"
)

func main() {
	// Open/create a storage.
	strg, err := storage.NewStorage()
	check(err)
	defer func() {
		check(strg.Close())
	}()

	// Read the data from the storage.
	users, err := user.Decode(strg.Reader())
	check(err)
	defer func() {
		check(strg.SaveSnapshot(&users))
	}()

	c := make(chan int)
	// Start a TCP server.
	go tcp.Server(c, strg, &users)
	// Start a TCP client.
	go tcp.Client(c)

	<-c
	<-c
	// Show the text user interface prompt.
	// tui.Prompt(os.Stdout, os.Stdin, strg, &users)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
