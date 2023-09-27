package storage

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	stdFileDir  = "datafiles"
	stdFileName = "test.database"
	dirPerms    = 0775 // rwxrwxr-x
	filePerms   = 0664 // rw-rw-r--
)

type Storage struct {
	file    *os.File
	updated bool
}

// Open opens the storage file.
func NewStorage() (storage *Storage, err error) {
	var fileName string
	if len(os.Args) < 2 {
		fmt.Println("No name for the database file. 'test.database' file is used.")
		fileName = stdFileName
	} else {
		fileName = os.Args[1]
		if strings.HasPrefix(fileName, stdFileDir+"/") {
			fileName = strings.TrimPrefix(fileName, stdFileDir+"/")
		}
	}
	if err = os.MkdirAll(stdFileDir, dirPerms); err != nil {
		return nil, err
	}
	fileName = filepath.Join(stdFileDir, fileName)

	storage = new(Storage)
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	storage.file, err = os.OpenFile(fileName, flags, filePerms)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) Reader() io.Reader {
	return bufio.NewReader(s.file)
}

func (s *Storage) Write(data []byte) (n int, err error) {
	writer := bufio.NewWriter(s.file)
	n, err = writer.Write(data)
	if err != nil {
		return 0, err
	}
	if err = writer.Flush(); err != nil {
		return 0, err
	}
	s.updated = true

	return n, nil
}

func (s *Storage) Close() {
	err := s.file.Close()
	if err != nil {
		log.Fatal(err)
	}
}
