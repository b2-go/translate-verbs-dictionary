package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	lang := os.Args[1]
	verb := os.Args[2]
	translateLang := os.Args[3]

	db, err := sql.Open("sqlite3", "./verbs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT langues, pers, plur, form FROM " + lang + " WHERE verb = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var langues string
	var pers int
	var plur int
	var form string
	err = stmt.QueryRow(verb).Scan(&langues, &pers, &plur, &form)
	if err != nil {
		log.Fatal(err)
	}

	translation := "NOT FOUND"
	switch translateLang {
	case "quechua":
		idx := strings.Index(langues, "QU=")
		if idx < 0 {
			log.Fatal("not translatable in ", translateLang)
		}
		verbInf := langues[idx+3:]
		translation = conjugueVerb(db, translateLang, verbInf, pers, plur, form)
	}

	log.Print(lang, " ", pers, plur, " ", form, " -> ", langues, " | ", translateLang, " = ", translation)
}

func conjugueVerb(db *sql.DB, lang string, verbInf string, pers, plur int, form string) string {

	switch lang {
	case "quechua":
		switch form {
		case "PR":
			form = "P"
		}
	case "francais":
		switch form {
		case "P":
			form = "PR"
		}
	}

	log.Print("CONJUGUE ", verbInf, " ", pers, plur, " ", form)
	stmt, err := db.Prepare("SELECT verb FROM " + lang + " WHERE verb_inf = ? AND pers = ? AND plur = ? AND form = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var verb string
	err = stmt.QueryRow(verbInf, pers, plur, form).Scan(&verb)
	if err == sql.ErrNoRows {
		return "NOT FOUND"
	}
	if err != nil {
		log.Fatal(err)
	}

	return verb
}
