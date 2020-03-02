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
		verbInf = strings.Trim(verbInf, `"`)
		translation = conjugueVerb(db, translateLang, verbInf, pers, plur, form)
	case "francais":
		idx := strings.Index(langues, "FR=")
		if idx < 0 {
			log.Fatal("not translatable in ", translateLang)
		}
		verbInf := langues[idx+3:]
		verbInf = strings.Trim(verbInf, `"`)
		translation = conjugueVerb(db, translateLang, verbInf, pers, plur, form)
	}

	log.Print(lang, " ", pers, plur, " ", form, " -> ", langues, " | ", translateLang, " = ", translation)
}

func conjugueVerb(db *sql.DB, lang string, verbInf string, pers, plur int, form string) string {

	// plur = que pour 1ere forme nous : p (français) <-> pex pluriel exclusive, pin pluriel inclusive (quechua)

	// aucun changement pour
	// - conditionnel
	// - futur
	switch lang {
	case "quechua":
		switch form {
		case "PR":
			form = "TI"
		case "I":
			form = "IP"
		case "SUBI":
			form = "S"
		case "PASS":
			form = "PS"
		case "GER":
			pers = 0 // un seul
			form = "G"
		case "GER1":
			pers = 0 // un seul
			form = "G"
		case "GER2":
			pers = 0 // un seul
			form = "G"
		case "PROG":
			// forme : "en train de {INF}"
		case "POT":
			// potentiel : "peut / capable de {INF}"
		}
		switch pers {

		}
	case "francais":
		switch form {
		case "TI":
			form = "PR"
		case "IP":
			form = "I"
		case "S":
			form = "SUBI"
		case "PS":
			form = "PASS"
		case "G":
			form = "GER"
			// "GER1"
			// "GER2"
			pers = -1 // chercher tous
		}
	}

	log.Print("CONJUGUE ", lang, " ", verbInf, " ", pers, plur, " ", form)
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
