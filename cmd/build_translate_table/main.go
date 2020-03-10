package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	translator "github.com/reagere/translate-verbs-dictionary"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := translator.NewDB()
	if err != nil {
		log.Panic(err)
	}
	tx := translator.NewService(db)

	errs, err := os.Create("errors.log")
	if err != nil {
		log.Fatal(err)
	}
	defer errs.Close()

	out, err := os.Create("output.log")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	csvout := csv.NewWriter(out)
	err = csvout.Write([]string{
		"langue origine",
		"verbe conjugue origine",
		"PersonneOrigine",
		"PlurielOrigine",
		"FormeOrigine",

		"langue traduite",
		"verbe conjugue traduit",
		"Personne Traduit",
		"Pluriel Traduit",
		"Forme Traduit",
	})
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	countError := 0
	for _, lang := range os.Args[1:] {
		list, err := translator.List(db, lang)
		if err != nil {
			log.Fatal(err)
		}
		// id, verb_inf, verb, itr, tr, langues, flex, flexOpts, pers, plur, form
		for i := range list {
			verb := list[i][2].(string)
			translateLang := "francais"
			if lang == "francais" {
				translateLang = "quechua"
			}
			translation, err := tx.Translate(lang, verb, translateLang)
			count++
			if err == nil {
				err = csvout.Write([]string{
					translation.LangueOrigine,
					translation.VerbeConjugueOrigine,
					strconv.Itoa(translation.PersonneOrigine),
					strconv.Itoa(translation.PlurielOrigine),
					translation.FormeOrigine,

					translation.LangueDestination,
					translation.VerbeConjugueTraduit,
					strconv.Itoa(translation.PersonneTraduit),
					strconv.Itoa(translation.PlurielTraduit),
					translation.FormeTraduit,
				})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				countError++
				errs.Write([]byte(fmt.Sprintf("%s %s -> %v\n", lang, verb, err)))
			}
			if count%100 == 0 {
				log.Printf("err : %d ; count : %d", countError, count)
			}
		}
	}
}
