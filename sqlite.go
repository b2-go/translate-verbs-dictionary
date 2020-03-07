package translator

import (
	"database/sql"
)

func NewDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./.data/verbs.db")
}
