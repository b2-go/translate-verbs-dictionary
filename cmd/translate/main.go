package main

import (
	"log"
	"os"

	translator "github.com/reagere/translate-verbs-dictionary"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	lang := os.Args[1]
	verb := os.Args[2]
	translateLang := os.Args[3]

	db, err := translator.NewDB()
	if err != nil {
		log.Panic(err)
	}

	tx := translator.NewService(db)
	tx.Translate(lang, verb, translateLang)
}
