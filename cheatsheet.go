package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func readInput(p string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p)
	s, _ := reader.ReadString('\n')
	return s
}

func addPrompts(inputs chan string) {

	com := strings.TrimSuffix(readInput("🧠 Enter command you would like to remember: "), "\n")
	inputs <- com

	lang := strings.TrimSuffix(readInput("️🛠  What language or tool is this for? "), "\n")
	inputs <- lang

	res := strings.TrimSuffix(readInput("⚡️ What does this command do? "), "\n")
	inputs <- res

}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Must specify command.")
		os.Exit(1)
	}

	action := os.Args[1]

	// open db
	db, err := sql.Open("sqlite3", "./cheatsheet.db")
	checkError(err)

	// create db
	create, err := db.Prepare("CREATE TABLE IF NOT EXISTS cheatsheet (command TEXT PRIMARY KEY, language TEXT, result TEXT)")
	checkError(err)
	create.Exec()

	// keep db open for remainder of function
	defer db.Close()

	if action == "ls" {
		rows, err := db.Query("SELECT * FROM cheatsheet")
		checkError(err)

		var command string
		var language string
		var result string

		fmt.Println("\nHere is your cheatsheet!📝")

		for rows.Next() {
			err = rows.Scan(&command, &language, &result)
			fmt.Printf("|  %-6s  |  %-6s  |  %-6s  |\n", command, language, result)

		}

		rows.Close()
	} else if action == "add" {
		inputs := make(chan string, 3)

		go addPrompts(inputs)

		insert, err := db.Prepare("INSERT INTO cheatsheet (command, language, result) values(?,?,?)")
		checkError(err)

		_, err = insert.Exec(<-inputs, <-inputs, <-inputs)
		checkError(err)

		fmt.Println("Successfully inserted into cheatsheet!🎉")

	} else {
		checkError(err)
	}

}
