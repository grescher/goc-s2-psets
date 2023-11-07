package cmd

import (
	"errors"
	"os"
	apiserv "practice/internal/api/server"
	"practice/internal/database/storage"
	"practice/internal/database/tcp"
	"practice/internal/database/tui"
	"practice/internal/database/user"
)

func Run() (err error) {
	// Open/create a storage.
	strg, err := storage.NewStorage()
	if err != nil {
		return err
	}
	// Read the data from the storage.
	users, err := user.Decode(strg.Reader())
	if err != nil {
		return err
	}
	defer func() {
		if checkErr := strg.SaveSnapshot(&users); checkErr != nil && err == nil {
			err = checkErr
		}
		if checkErr := strg.Close(); checkErr != nil && err == nil {
			err = checkErr
		}
	}()

	c := make(chan int)
	// Start a TCP server.
	go tcp.Server(c, strg, &users)
	go apiserv.Run(c)

	checkErr := tui.Prompt(os.Stdout, os.Stdin, strg, &users)
	if !errors.Is(checkErr, tui.ErrEndOfSession) {
		err = checkErr
	}
	close(c)

	return err
}
