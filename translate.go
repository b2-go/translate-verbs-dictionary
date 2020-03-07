package translator

import (
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	db   *sql.DB
	conj *Conjugue
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db:   db,
		conj: NewConjugue(db),
	}
}

type Translation struct {
	LangueOrigine        string
	VerbeConjugueOrigine string
	// donnée
	PersonneOrigine int
	PlurielOrigine  int
	FormeOrigine    string
	// results
	LangueDestination    string
	VerbeConjugueTraduit string
	// details
	PersonneTraduit int
	PlurielTraduit  int
	FormeTraduit    string
}

func (s *Service) Translate(lang, verb, translateLang string) (*Translation, error) {

	stmt, err := s.db.Prepare("SELECT langues, pers, plur, form FROM " + lang + " WHERE verb = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var langues string
	var pers int
	var plur int
	var form string
	err = stmt.QueryRow(verb).Scan(&langues, &pers, &plur, &form)
	if err != nil {
		return nil, err
	}

	prefixLang := ""
	switch translateLang {
	case "quechua":
		prefixLang = "QU="
	case "francais":
		prefixLang = "FR="
	default:
		return nil, fmt.Errorf("langue non fournie %s", translateLang)
	}

	idx := strings.Index(langues, prefixLang)
	if idx < 0 {
		return nil, fmt.Errorf("not translatable in %s", translateLang)
	}
	verbInf := langues[idx+3:]
	verbInf = strings.Trim(verbInf, `"`)

	formTx, persTx := s.convertForm(lang, form, pers)

	translation, err := s.conj.ConjugueVerb(translateLang, verbInf, persTx, plur, formTx)
	if err != nil {
		return nil, err
	}

	return &Translation{
		LangueOrigine:        lang,
		VerbeConjugueOrigine: verb,
		PersonneOrigine:      pers,
		PlurielOrigine:       plur,
		FormeOrigine:         form,
		//
		LangueDestination:    translateLang,
		VerbeConjugueTraduit: translation,
		PersonneTraduit:      persTx,
		PlurielTraduit:       plur,
		FormeTraduit:         formTx,
	}, nil
}

func (s *Service) convertForm(lang, form string, pers int) (string, int) {

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
	return form, pers
}
