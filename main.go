package main

func main() {
	users := Users()
	headers := []string{"Name", "Age", "Active", "Mass", "Books"}

	PrintData(users, headers)
}
