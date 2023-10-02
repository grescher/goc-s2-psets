package storage

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"practice/internal/user"
)

const (
	stdFileDir  = "datafiles"
	stdFileName = "test.database"
	dirPerms    = 0775 // rwxrwxr-x
	filePerms   = 0664 // rw-rw-r--
)

type Storage struct {
	file    *os.File
	path    string
	updated bool
}

// Open opens the storage file.
func NewStorage() (strg *Storage, err error) {
	var fileDir, fileName string
	if len(os.Args) < 2 {
		fmt.Println("No name for the database file. 'test.database' file is used.")
		fileDir, fileName = stdFileDir, stdFileName
	} else {
		fileDir, fileName = filepath.Split(os.Args[1])
	}

	if len(fileDir) > 0 {
		err = os.MkdirAll(fileDir, dirPerms)
		if err != nil {
			return nil, err
		}
	}
	filePath := filepath.Join(fileDir, fileName)

	strg = new(Storage)
	strg.path = filePath
	if err = strg.Open(); err != nil {
		return nil, err
	}

	return strg, nil
}

func (s *Storage) Open() (err error) {
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	s.file, err = os.OpenFile(s.path, flags, fs.FileMode(filePerms))
	return err
}

func (s *Storage) Reader() io.Reader {
	return s.file
}

func (s *Storage) Writer() io.Writer {
	return s.file
}

func (s *Storage) Sync() error {
	return s.file.Sync()
}

func (s *Storage) SaveSnapshot(users *[]user.User) (err error) {
	// Create a temporary storage file.
	tmpDir, tmpFileName := filepath.Split(s.path)
	tmpFile, err := os.CreateTemp(tmpDir, tmpFileName)
	if err != nil {
		return err
	}
	tmpFilePath := tmpFile.Name()

	// Encode data and save it to the temp file.
	if err = user.Encode(tmpFile, *users); err != nil {
		return err
	}
	if err = tmpFile.Sync(); err != nil {
		return err
	}

	// Close both files: either the temp or the storage.
	if err = tmpFile.Close(); err != nil {
		return err
	}
	if err = s.Close(); err != nil {
		return err
	}

	// Rename (move) temporary into the main storage file.
	if err = os.Rename(tmpFilePath, s.path); err != nil {
		return err
	}
	if err = os.Chmod(s.path, fs.FileMode(filePerms)); err != nil {
		return err
	}
	if err = s.Open(); err != nil {
		return err
	}
	return nil
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

func (s *Storage) Close() error {
	return s.file.Close()
}
