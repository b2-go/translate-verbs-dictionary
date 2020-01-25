package main

import (
	"log"
	"os"

	translator "github.com/reagere/translate-verbs-dictionary"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := translator.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	for _, lang := range os.Args[1:] {
		list, err := translator.List(db, lang)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("=== %s ===", lang)
		log.Print("HEAD: id, verb_inf, verb, itr, tr, langues, flex, flexOpts, pers, plur, form")
		for i := range list {
			log.Print(list[i])
		}
	}
}
