// tui implements a text-based user interface.
package tui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"practice/internal/database/storage"
	"practice/internal/database/tui/table"
	"practice/internal/database/user"
	"strconv"
	"strings"
)

var (
	ErrEndOfSession = errors.New("end of session")
	ErrUserNotFound = errors.New("user is not found")
)

func Prompt(w io.Writer, r io.Reader, strg *storage.Storage, users *[]user.User) error {
	fmt.Fprintln(w, "Enter \"help\" for usage hints.")

	for {
		fmt.Fprintf(w, "%s > ", strg.Name())

		var in string
		fmt.Fscanf(r, "%s", &in)

		switch strings.ToUpper(in) {
		case "ADD":
			if err := addUser(w, r, strg, users); err != nil {
				log.Println("failed to add user:", err)
			}
		case "REMOVE":
			err := rmUser(w, r, strg, users)
			switch {
			case err == ErrUserNotFound:
				fmt.Fprintln(w, err)
			case err != nil:
				log.Println("failed to remove user:", err)
			default:
				fmt.Fprintln(w, "User deleted")
			}
		case "SHOW":
			show(w, *users)
		case "HELP":
			printHelp(w)
		case "QUIT", "EXIT":
			return ErrEndOfSession
		default:
			fmt.Fprintf(w, "Unknown operator %q. Enter \"help\" for usage hints.\n", in)
		}
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintln(
		w,
		`add         Add user to the database
help        Print this help
quit, exit  Exit the program
remove      Remove the user from the database
show        Print contents of the database`,
	)
}

// addUser adds a new user to the slice of users and writes them to the storage.
func addUser(w io.Writer, r io.Reader, strg *storage.Storage, users *[]user.User) error {
	// Check if there is space for a new user.
	if len(*users) > user.MaxNumOfUsers {
		return errors.New("no free slots for a new user")
	}

	// Read the new user's data.
	rb := bufio.NewReader(r)
	// - name:
	name, err := promptUserName(w, rb)
	if err != nil {
		return err
	}
	if name == "" {
		return errors.New("no name is entered")
	}
	// - age:
	age, err := promptUserAge(w, rb)
	if err != nil {
		return err
	}
	// - active index/status:
	activeIndex, err := promptUserActiveStatus(w, rb, *users)
	if err != nil {
		return err
	}
	// - mass:
	mass, err := promptUserMass(w, rb)
	if err != nil {
		return err
	}
	// - books:
	var books []string
	if err = promptUserBooks(w, rb, &books); err != nil {
		return err
	}

	// Add a new user to the users.
	newUser := user.User{
		Name:        name,
		Age:         age,
		ActiveIndex: activeIndex,
		Mass:        mass,
		Books:       books,
	}
	*users = append(*users, newUser)

	// Save a new user to the storage.
	if err = user.EncodeUser(strg.Writer(), newUser); err != nil {
		return err
	}
	if err = strg.Sync(); err != nil {
		return err
	}

	return nil
}

// promptUserName prompts for a name of a new user.
func promptUserName(w io.Writer, r *bufio.Reader) (input string, err error) {
	fmt.Fprint(w, "Enter name: ")

	input, err = r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("couldn't read name: %v", err)
	}

	return strings.TrimSpace(input), nil
}

// promptUserAge prompts for an age of a new user.
func promptUserAge(w io.Writer, r *bufio.Reader) (uint8, error) {
	fmt.Fprint(w, "Enter age: ")

	input, err := r.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("couldn't read age: %v", err)
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return 0, nil
	}
	age, err := strconv.ParseUint(input, 10, 8)
	if err != nil {
		return 0, err
	}

	return uint8(age), nil
}

// promptUserActiveStatus prompts if a new user is active and generates the activeStatus index.
func promptUserActiveStatus(w io.Writer, r *bufio.Reader, users []user.User) (activeStatus uint8, err error) {
	fmt.Fprint(w, "Is the user is active now? [yes/no]: ")

	input, err := r.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("couldn't read active status: %v", err)
	}
	input = strings.TrimSpace(input)

	switch strings.ToUpper(input) {
	case "YES", "Y":
		activeStatus = 1
	case "NO", "N", "":
		activeStatus = 0
	default:
		fmt.Fprint(w, "Please, provide with [yes/no], [YyNn].")
		return promptUserActiveStatus(w, r, users)
	}
	// TODO: improve the activeStatus index generation to avoid duplicates.
	activeStatus <<= len(users)

	return activeStatus, nil
}

// promptUserMass prompts for a mass of a new user.
func promptUserMass(w io.Writer, r *bufio.Reader) (mass float64, err error) {
	fmt.Fprint(w, "Enter the user's mass: ")

	input, err := r.ReadString('\n')
	if err != nil {
		return 0.0, fmt.Errorf("couldn't read mass: %v", err)
	}
	input = strings.TrimSpace(input)

	mass, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return 0.0, fmt.Errorf("couldn't read mass: %v", err)
	}
	if mass < 0.0 {
		mass = 0.0
	}

	return user.VerifyMass(mass), nil
}

// promptUserBooks prompts for a list of books a new user has read.
func promptUserBooks(w io.Writer, r *bufio.Reader, books *[]string) error {
	fmt.Fprint(w, "Enter a name of book: ")

	input, err := r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("couldn't read book name: %v", err)
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return nil //-> exit point.
	}

	*books = append(*books, input)
	return promptUserBooks(w, r, books)
}

// rmUser searches for a user by name, and if it finds them, removes them from the slice
// of users; after that, the snapshot of the slice of users is saved in the storage.
func rmUser(w io.Writer, r io.Reader, strg *storage.Storage, users *[]user.User) (err error) {
	reader := bufio.NewReader(r)

	// Find the user by name. Determine the user's index.
	fmt.Fprint(w, "Enter the name of user you want to remove: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("couldn't read name: %v", err)
	}
	name := strings.TrimSpace(input)

	i, ok := user.Slice(*users).FindName(name)
	if !ok {
		return ErrUserNotFound
	}

	// Remove the user from the users slice.
	(*users)[i] = (*users)[len(*users)-1]
	*users = (*users)[:len(*users)-1]

	// Save the snapshot.
	if err = strg.SaveSnapshot(users); err != nil {
		return err
	}

	return nil
}

func show(w io.Writer, users []user.User) {
	table.PrintData(w, user.Slice(users), user.Headers)
	fmt.Fprintln(w, "Number of active users:", user.Slice(users).NumOfActiveUsers())
}
