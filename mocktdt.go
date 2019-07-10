// Command mocktdt is a sample program to illustrate the table-driven
// mock database unittesting.package mocktdt.
// (c) rusq for Gophership conference in Moscow, 20/07/2019
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/pressly/goose"
	"github.com/rusq/mocktdt/db/repository"
	"github.com/rusq/mocktdt/types"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbDriver      = "sqlite3"
	dbConnStr     = "production.db"
	migrationsDir = "db/migrations"
)

var db *sql.DB // global db handle

func init() {

	var err error
	db, err = sql.Open(dbDriver, dbConnStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// run available migrations
	if err := goose.SetDialect(dbDriver); err != nil {
		panic(err)
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		panic(err)
	}
}

func main() {
	if err := doSomeStuff(); err != nil {
		log.Fatal(err)
	}

	log.Print("we should not get here!")
}

func doSomeStuff() error {
	repo := repository.NewPersonaRepository()

	john := types.Persona{ID: 1, Name: "Johnnie Walker", DOB: mustPtrTime("1865-01-01 00:00:00")}

	if err := repo.Delete(context.TODO(), db, john.ID); err != nil {
		return err
	}
	if err := repo.Insert(context.TODO(), db, &john); err != nil {
		return err
	}
	if exists, err := repo.Exists(context.TODO(), db, john.ID); err != nil {
		return err
	} else {
		if !exists {
			return errors.New("John, why you no exist?!1")
		}
	}

	john.Name = "Jack Daniel"
	john.DOB = mustPtrTime("1875-01-01 00:00:00")
	if err := repo.Update(context.TODO(), db, &john); err != nil {
		return err
	}

	return nil
}

// mustPtrTime returns a time pointer from date in "2006-01-02 15:04:05" fmt
func mustPtrTime(s string) *time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return &t
}
