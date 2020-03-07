package translator

import (
	"database/sql"
	"log"
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

func (s *Service) Translate(lang, verb, translateLang string) {

	stmt, err := s.db.Prepare("SELECT langues, pers, plur, form FROM " + lang + " WHERE verb = ?")
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
		translation = s.conj.ConjugueVerb(translateLang, verbInf, pers, plur, form)
	case "francais":
		idx := strings.Index(langues, "FR=")
		if idx < 0 {
			log.Fatal("not translatable in ", translateLang)
		}
		verbInf := langues[idx+3:]
		verbInf = strings.Trim(verbInf, `"`)
		translation = s.conj.ConjugueVerb(translateLang, verbInf, pers, plur, form)
	}

	log.Print(lang, " ", pers, plur, " ", form, " -> ", langues, " | ", translateLang, " = ", translation)
}
