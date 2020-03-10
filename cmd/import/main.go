package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	translator "github.com/reagere/translate-verbs-dictionary"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	lang := os.Args[1]
	file := os.Args[2]

	db, err := translator.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
CREATE TABLE %s (id INTEGER PRIMARY KEY AUTOINCREMENT, verb_inf TEXT, verb TEXT, itr NUMERIC, tr NUMERIC, langues TEXT, flex TEXT, flexOpts TEXT, pers INTEGER, plur INTEGER, form TEXT);
CREATE INDEX idx_verb_%s ON %s (verb);
CREATE INDEX idx_form_%s ON %s (pers, plur, form);
	`
	_, err = db.Exec(fmt.Sprintf(sqlStmt, lang, lang, lang, lang, lang))
	if err != nil {
		log.Fatal(err)
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

	stmt, err := tx.Prepare("INSERT INTO " + lang + " (verb_inf, verb, itr, tr, langues, flex, flexOpts, pers, plur, form) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for {
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
		pers := 0
		plur := 0 // indefini, 1: singulier, 2: pluriel
		form := ""
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
					opt = strings.Replace(opt, `"`, "", -1)
					opt = strings.Trim(opt, ` `)
					langues += opt
				} else if strings.HasPrefix(opt, "FLX=") {
					flex = opt[4:]
					flexOpts = options[i+1:]
					for _, flexOpt := range flexOpts {
						switch flexOpt {
						case "1":
							pers = 1
						case "2":
							pers = 2
						case "3":
							pers = 3
						case "s":
							plur = 1
						case "p":
							plur = 2
						case "pin":
							plur = 3
						case "pex":
							plur = 4
							// francais
						case "PR":
							form = flexOpt
						case "F":
							form = flexOpt
						case "G":
							form = flexOpt
						case "C":
							form = flexOpt
						case "S":
							form = flexOpt
						case "IP":
							form = flexOpt
						case "I":
							form = flexOpt
						case "PP":
							form = flexOpt
						case "INF":
							form = flexOpt
							// quechua
						case "FA":
							form = flexOpt
						case "FP":
							form = flexOpt
						case "GER":
							form = flexOpt
						case "GER2":
							form = flexOpt
						case "OBL":
							form = flexOpt
						case "PRES":
							form = flexOpt
						case "PASS":
							form = flexOpt
						case "PASS1":
							form = flexOpt
						case "PASS2":
							form = flexOpt
						case "PPA":
							form = flexOpt
						case "PPA2":
							form = flexOpt
						case "PAPT":
							form = flexOpt
						case "PPAT":
							form = flexOpt
						case "PPI":
							form = flexOpt
						case "RQUF":
							form = flexOpt
						case "RS1":
							form = flexOpt
						case "SUBI":
							form = flexOpt
						case "TI":
							form = flexOpt
						}
					}
				}
			}
		}
		_, err = stmt.Exec(l[1], l[0], itr, tr, langues, flex, strings.Join(flexOpts, "+"), pers, plur, form)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
