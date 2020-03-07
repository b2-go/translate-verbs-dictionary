package translator

import (
	"database/sql"
	"fmt"
)

type Conjugue struct {
	db *sql.DB
}

func NewConjugue(db *sql.DB) *Conjugue {
	return &Conjugue{
		db: db,
	}
}

func (c *Conjugue) ConjugueVerb(lang string, verbInf string, pers, plur int, form string) (string, error) {
	stmt, err := c.db.Prepare("SELECT verb FROM " + lang + " WHERE verb_inf = ? AND pers = ? AND plur = ? AND form = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var verb string
	err = stmt.QueryRow(verbInf, pers, plur, form).Scan(&verb)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("not found")
	}
	if err != nil {
		return "", err
	}

	return verb, nil
}
