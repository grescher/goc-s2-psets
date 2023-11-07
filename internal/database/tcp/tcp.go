package tcp

import (
	"io"
	"log"
	"net"
	"os"
	"practice/internal/database/storage"
	"practice/internal/database/tui"
	"practice/internal/database/user"
)

func Server(c chan int, strg *storage.Storage, users *[]user.User) {
	defer func() {
		c <- 0
	}()

	// Create listener.
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("tcp.Server: failed to create a listener: ", err)
	}
	defer listener.Close()

	// Listen for a new connection.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("tcp.Server: failed to accept a connection:", err)
			continue
		}

		err = handleConn(conn, strg, users)
		if err != nil && err != tui.ErrEndOfSession {
			log.Println("tcp.Server: handling connection:", err)
		}
		conn.Close()
	}
}

func handleConn(conn net.Conn, strg *storage.Storage, users *[]user.User) (err error) {
	if err := tui.Prompt(conn, conn, strg, users); err != nil {
		return err
	}
	return nil
}

func Client(c chan int) {
	defer func() {
		c <- 0
	}()

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("tcp.Client: failed connect to server: ", err)
	}
	defer conn.Close()

	go func() {
		if _, err = io.Copy(os.Stdout, conn); err != nil {
			log.Println("tcp.Client: failed to redirect from connection to Stdout:", err)
		}
	}()

	if _, err = io.Copy(conn, os.Stdin); err != nil {
		log.Println("tcp.Client: failed to redirect from Stdin to connection:", err)
	}
}
