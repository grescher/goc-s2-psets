// tui implements a text-based user interface.
package tui

import (
	"bufio"
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
			addUser(strg, users)
		case "REMOVE":
			rmUser(strg, users)
		case "SHOW":
			show(*users)
		case "HELP":
			printHelp()
		case "QUIT":
			return
		default:
			fmt.Printf("Unknown operator %q. Enter \".help\" for usage hints.\n", in)
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
func addUser(strg *storage.Storage, users *[]user.User) {
	if len(*users) > user.MaxNumOfUsers {
		fmt.Println("There's no free slots for a new user.")
		return
	}
	reader := bufio.NewReader(os.Stdin)

	// Read the new user's data.
	name := promptUserName(reader)
	if name == "" {
		return
	}
	age := promptUserAge(reader)
	activeIndex := promptUserActiveStatus(reader, *users)
	mass := promptUserMass(reader)

	var books []string
	promptUserBooks(reader, &books)

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
	user.EncodeUser(strg.Writer(), newUser)
	strg.Sync()
}

// promptUserName prompts for a name of a new user.
func promptUserName(reader *bufio.Reader) string {
	fmt.Print("Enter name: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Couldn't read the user's name:", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// promptUserAge prompts for an age of a new user.
func promptUserAge(reader *bufio.Reader) uint8 {
	fmt.Print("Enter age: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Couldn't read the user's age:", err)
		return 0
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return 0
	}
	age, err := strconv.ParseUint(input, 10, 8)
	if err != nil {
		fmt.Println("Unknown type of the user's age. Try again.")
		return promptUserAge(reader)
	}
	return uint8(age)
}

// promptUserActiveStatus prompts if a new user is active and generates the activeStatus index.
func promptUserActiveStatus(reader *bufio.Reader, users []user.User) (activeStatus uint8) {
	fmt.Print("Is the user is active now? [yes/no]: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Couldn't read the user's active status:", err)
		input = ""
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

	return activeStatus
}

// promptUserMass prompts for a mass of a new user.
func promptUserMass(reader *bufio.Reader) (mass float64) {
	fmt.Print("Enter the user's mass: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Couldn't read the user's mass:", err)
		return 0.0
	}
	input = strings.TrimSpace(input)
	mass, err = strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Unknown type of the mass. Please, try again with a number.")
		return promptUserMass(reader)
	}
	if mass < 0.0 {
		mass = 0.0
	}
	return user.VerifyMass(mass)
}

// promptUserBooks prompts for a list of books a new user has read.
func promptUserBooks(reader *bufio.Reader, books *[]string) {
	fmt.Print("Enter a name of book: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Couldn't read the input:", err)
		return
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	*books = append(*books, input)
	promptUserBooks(reader, books)
}

// rmUser searches for a user by name, and if it finds them, removes them from the slice
// of users; after that, the snapshot of the slice of users is saved in the storage.
func rmUser(strg *storage.Storage, users *[]user.User) {
	reader := bufio.NewReader(os.Stdin)

	// Find the user by name. Determine the user's index.
	fmt.Print("Enter the name of user you want to remove: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Couldn't remove the user:", err)
		return
	}
	name := strings.TrimSpace(input)

	i, ok := user.Slice(*users).FindName(name)
	if !ok {
		fmt.Printf("Couldn't find the user %q\n", name)
		return
	}

	// Remove the user from the users slice.
	(*users)[i] = (*users)[len(*users)-1]
	*users = (*users)[:len(*users)-1]

	// Save the snapshot.
	strg.SaveSnapshot(users)
}

func show(users []user.User) {
	table.PrintData(user.Slice(users), user.Headers)
	fmt.Println("Number of active users:", user.Slice(users).NumOfActiveUsers())
}
