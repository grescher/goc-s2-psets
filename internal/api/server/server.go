package apiserv

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	addr   = "localhost:8080"
	dbAddr = "localhost:8000"
)

func Run(c chan int) (err error) {
	defer func() {
		c <- 0
	}()
	mux := http.NewServeMux()

	mux.HandleFunc("/users", usersHandler)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	err = server.ListenAndServe()
	return err
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("tcp", dbAddr)
	if err != nil {
		log.Println("usersHandler:", err)
		return
	}
	defer conn.Close()

	if _, err = fmt.Fprint(conn, "show\n"); err != nil {
		log.Println("usersHandler:", err)
		return
	}

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("usersHandler: reader:", err)
			return
		}
		fmt.Fprint(w, string(buf[:n]))
	}
}
