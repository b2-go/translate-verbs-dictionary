package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	lang := os.Args[1]
	verb := os.Args[2]

	db, err := sql.Open("sqlite3", "./verbs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT langues, flex, flexOpts FROM " + lang + " WHERE verb = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var langues string
	var flex string
	var flexOpts string
	err = stmt.QueryRow(verb).Scan(&langues, &flex, &flexOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(langues, flex, flexOpts)
}
