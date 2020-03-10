package translator

import (
	"database/sql"
	"fmt"
)

func List(db *sql.DB, lang string) ([][]interface{}, error) {
	rows, err := db.Query("SELECT id, verb_inf, verb, itr, tr, langues, flex, flexOpts, pers, plur, form FROM " + lang)
	if err != nil {
		return nil, fmt.Errorf("list query: %w", err)
	}
	defer rows.Close()

	result := make([][]interface{}, 0)

	for rows.Next() {
		var id int
		var verbInf string
		var verb string
		var itr bool
		var tr bool
		var langues string
		var flex string
		var flexOpts string
		var pers int
		var plur int
		var form string
		err = rows.Scan(&id, &verbInf, &verb, &itr, &tr, &langues, &flex, &flexOpts, &pers, &plur, &form)
		if err != nil {
			return nil, fmt.Errorf("list scan: %w", err)
		}
		result = append(result, []interface{}{id, verbInf, verb, itr, tr, langues, flex, flexOpts, pers, plur, form})
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("list errors: %w", err)
	}

	return result, nil
}
