package translator

import (
	"database/sql"
	"log"
)

type Conjugue struct {
	db *sql.DB
}

func NewConjugue(db *sql.DB) *Conjugue {
	return &Conjugue{
		db: db,
	}
}

func (c *Conjugue) ConjugueVerb(lang string, verbInf string, pers, plur int, form string) string {

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
	stmt, err := c.db.Prepare("SELECT verb FROM " + lang + " WHERE verb_inf = ? AND pers = ? AND plur = ? AND form = ?")
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
