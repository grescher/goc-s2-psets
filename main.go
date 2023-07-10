package main

import "fmt"

func main() {
	users := UserSlice(Users())
	PrintData(users, UserHeaders)
	fmt.Println("Books Stats:")

}

func PrintData(data TablePrinter, headers []string) {
	table := data.NewTable(headers)
	table.Print()
}
