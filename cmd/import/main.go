package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	lang := os.Args[1]
	file := os.Args[2]

	db, err := sql.Open("sqlite3", "./verbs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := "CREATE TABLE " + lang + " (id INTEGER not null primary key, verb_inf TEXT, verb TEXT, itr NUMERIC, tr NUMERIC, langues TEXT, flex TEXT, flexOpts TEXT)"
	_, _ = db.Exec(sqlStmt)

	sqlStmt = "DELETE FROM " + lang
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s", err, sqlStmt)
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	csvr := csv.NewReader(f)
	csvr.Comma = ','
	csvr.LazyQuotes = true

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO " + lang + " (id, verb_inf, verb, itr, tr, langues, flex, flexOpts) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; ; i++ {
		l, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		itr := false
		tr := false
		langues := ""
		flex := ""
		flexOpts := []string{}
		options := strings.Split(l[2], "+")
		for i := range options {
			opt := options[i]
			switch opt {
			case "ITR":
				itr = true
			case "TR":
				tr = true
			default:
				if strings.HasPrefix(opt, "FR=") || strings.HasPrefix(opt, "QU=") {
					if langues != "" {
						langues += ","
					}
					langues += opt
				} else if strings.HasPrefix(opt, "FLEX=") {
					flex = opt[5:]
					flexOpts = options[i+1:]
				}
			}
		}
		_, err = stmt.Exec(i, l[1], l[0], itr, tr, langues, flex, strings.Join(flexOpts, "+"))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, verb_inf, verb, itr, tr, langues, flex, flexOpts FROM " + lang)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var verbInf string
		var verb string
		var itr bool
		var tr bool
		var langues string
		var flex string
		var flexOpts string
		err = rows.Scan(&id, &verbInf, &verb, &itr, &tr, &langues, &flex, &flexOpts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, verbInf, verb, itr, tr, langues, flex, flexOpts)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
