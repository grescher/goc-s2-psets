// tui implements a text-based user interface.
package tui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"practice/internal/storage"
	"practice/internal/table"
	"practice/internal/user"
	"strconv"
	"strings"
)

func Prompt(strg *storage.Storage, users *[]user.User) {
	fmt.Println("Enter \"help\" for usage hints.")

	for {
		fmt.Printf("%s > ", strg.Name())

		var in string
		fmt.Scanf("%s", &in)

		switch strings.ToUpper(in) {
		case "ADD":
			if err := addUser(strg, users); err != nil {
				fmt.Println("failed to add user:", err)
			}
		case "REMOVE":
			if err := rmUser(strg, users); err != nil {
				fmt.Println("failed to remove user:", err)
				continue
			}
			fmt.Println("User deleted")
		case "SHOW":
			show(*users)
		case "HELP":
			printHelp()
		case "QUIT":
			return
		default:
			fmt.Printf("Unknown operator %q. Enter \"help\" for usage hints.\n", in)
		}
	}
}

func printHelp() {
	fmt.Println(
		`add     Adds user to the database
help    Show help
quit    Exit this program
remove  Removes the user from the database
show    Prints the contents of the table`,
	)
}

// addUser adds a new user to the slice of users and writes them to the storage.
func addUser(strg *storage.Storage, users *[]user.User) (err error) {
	if len(*users) > user.MaxNumOfUsers {
		return errors.New("no free slots for a new user")
	}
	reader := bufio.NewReader(os.Stdin)

	// Read the new user's data.

	// - name:
	name, err := promptUserName(reader)
	if err != nil {
		return err
	}
	if name == "" {
		return errors.New("no name is entered")
	}
	// - age:
	age, err := promptUserAge(reader)
	if err != nil {
		return err
	}
	// - active index/status:
	activeIndex, err := promptUserActiveStatus(reader, *users)
	if err != nil {
		return err
	}
	// - mass:
	mass, err := promptUserMass(reader)
	if err != nil {
		return err
	}
	// - books:
	var books []string
	if err = promptUserBooks(reader, &books); err != nil {
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
func promptUserName(reader *bufio.Reader) (input string, err error) {
	fmt.Print("Enter name: ")

	input, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("couldn't read name: %v", err)
	}

	return strings.TrimSpace(input), nil
}

// promptUserAge prompts for an age of a new user.
func promptUserAge(reader *bufio.Reader) (uint8, error) {
	fmt.Print("Enter age: ")

	input, err := reader.ReadString('\n')
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
func promptUserActiveStatus(reader *bufio.Reader, users []user.User) (activeStatus uint8, err error) {
	fmt.Print("Is the user is active now? [yes/no]: ")

	input, err := reader.ReadString('\n')
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
		fmt.Println("Please, provide with [yes/no], [YyNn].")
		return promptUserActiveStatus(reader, users)
	}
	// TODO: improve the activeStatus index generation to avoid duplicates.
	activeStatus <<= len(users)

	return activeStatus, nil
}

// promptUserMass prompts for a mass of a new user.
func promptUserMass(reader *bufio.Reader) (mass float64, err error) {
	fmt.Print("Enter the user's mass: ")

	input, err := reader.ReadString('\n')
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
func promptUserBooks(reader *bufio.Reader, books *[]string) error {
	fmt.Print("Enter a name of book: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("couldn't read book name: %v", err)
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return nil //-> exit point.
	}

	*books = append(*books, input)
	return promptUserBooks(reader, books)
}

// rmUser searches for a user by name, and if it finds them, removes them from the slice
// of users; after that, the snapshot of the slice of users is saved in the storage.
func rmUser(strg *storage.Storage, users *[]user.User) (err error) {
	reader := bufio.NewReader(os.Stdin)

	// Find the user by name. Determine the user's index.
	fmt.Print("Enter the name of user you want to remove: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("couldn't read name: %v", err)
	}
	name := strings.TrimSpace(input)

	i, ok := user.Slice(*users).FindName(name)
	if !ok {
		return fmt.Errorf("couldn't find the user %q", name)
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

func show(users []user.User) {
	table.PrintData(user.Slice(users), user.Headers)
	fmt.Println("Number of active users:", user.Slice(users).NumOfActiveUsers())
}
