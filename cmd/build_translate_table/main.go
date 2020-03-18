package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	translator "github.com/reagere/translate-verbs-dictionary"

	"golang.org/x/sync/errgroup"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := translator.NewDB()
	if err != nil {
		log.Panic(err)
	}
	tx := translator.NewService(db)

	errs, err := os.Create("errors.log")
	if err != nil {
		log.Fatal(err)
	}
	defer errs.Close()
	var errsLock sync.Mutex

	out, err := os.Create("output.log")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	var outLock sync.Mutex

	csvout := csv.NewWriter(out)
	defer csvout.Flush()
	err = csvout.Write([]string{
		"langue origine",
		"verbe conjugue origine",
		"PersonneOrigine",
		"PlurielOrigine",
		"FormeOrigine",

		"langue traduite",
		"verbe conjugue traduit",
		"Personne Traduit",
		"Pluriel Traduit",
		"Forme Traduit",
	})
	if err != nil {
		log.Fatal(err)
	}

	g, ctx := errgroup.WithContext(context.Background())

	_ = ctx
	for i := range os.Args[1:] {
		lang := os.Args[1:][i]
		translateLang := "francais"
		if lang == "francais" {
			translateLang = "quechua"
		}
		g.Go(func() error {
			list, err := translator.List(db, lang)
			if err != nil {
				return err
			}

			b := &Builder{
				lang:          lang,
				translateLang: translateLang,
				errsLock:      &errsLock,
				errs:          errs,
				outLock:       &outLock,
				csvout:        csvout,
				tx:            tx,
				list:          list,
			}

			size := len(list)
			splitSize := size / 5
			offset := 0
			for split := 0; split < 4; split++ {
				g.Go(b.run(offset, offset+splitSize))
				offset += split
			}
			g.Go(b.run(offset, size))
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

type Builder struct {
	lang          string
	translateLang string
	tx            *translator.Service
	errsLock      *sync.Mutex
	errs          io.Writer
	outLock       *sync.Mutex
	csvout        *csv.Writer
	list          [][]interface{}
}

func (b *Builder) run(offset, limit int) func() error {
	return func() error {
		// id, verb_inf, verb, itr, tr, langues, flex, flexOpts, pers, plur, form
		for i := offset; i < limit; i++ {
			verb, ok := b.list[i][2].(string)
			if !ok {
				b.errsLock.Lock()
				_, err := b.errs.Write([]byte("not a string"))
				b.errsLock.Unlock()
				if err != nil {
					return err
				}
				continue
			}
			translation, err := b.tx.Translate(b.lang, verb, b.translateLang)
			if err == nil {
				b.outLock.Lock()
				err = b.csvout.Write([]string{
					translation.LangueOrigine,
					translation.VerbeConjugueOrigine,
					strconv.Itoa(translation.PersonneOrigine),
					strconv.Itoa(translation.PlurielOrigine),
					translation.FormeOrigine,

					translation.LangueDestination,
					translation.VerbeConjugueTraduit,
					strconv.Itoa(translation.PersonneTraduit),
					strconv.Itoa(translation.PlurielTraduit),
					translation.FormeTraduit,
				})
				b.outLock.Unlock()
				if err != nil {
					return err
				}
			} else {
				b.errsLock.Lock()
				_, err = b.errs.Write([]byte(fmt.Sprintf("%s %s -> %v\n", b.lang, verb, err)))
				b.errsLock.Unlock()
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}
