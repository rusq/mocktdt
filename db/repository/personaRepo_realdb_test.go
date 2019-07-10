package repository

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"

	"github.com/rusq/mocktdt/types"
)

const (
	testDbDrv  = "sqlite3"
	testDbConn = "test.db"
)

var migrationsDir = filepath.Join("..", "migrations")

var peeps = []types.Persona{
	{ID: 1, Name: "Daniil Kharms", DOB: mustPtrTime("1990-08-15 00:00:00")},
	{ID: 2, Name: "Jack London", DOB: mustPtrTime("1989-07-14 00:00:00")},
	{ID: 3, Name: "Isaak Asimov", DOB: mustPtrTime("1988-06-13 00:00:00")},
	{ID: 4, Name: "Philip K. Dick", DOB: mustPtrTime("1987-05-12 00:00:00")},
}

// check panics on error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func initTestDB(drv, conn string) *sql.DB {
	db, err := sql.Open(drv, conn)
	check(err)
	check(db.Ping())
	check(goose.SetDialect(drv))
	check(goose.Up(db, migrationsDir))
	return db
}

func deinitSqlite3(dbfile string) {
	os.Remove(dbfile)
}

// mustPtrTime returns a time pointer from date in "2006-01-02 15:04:05" fmt
func mustPtrTime(s string) *time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	check(err)
	return &t
}

func TestPersonaRepository_Insert_REAL(t *testing.T) {
	db := initTestDB(testDbDrv, testDbConn)
	// cleanup
	defer db.Close()
	defer deinitSqlite3(testDbConn)

	repo := NewPersonaRepository()

	for _, persona := range peeps {
		if err := repo.Insert(context.Background(), db, &persona); err != nil {
			t.Errorf("failed to insert: %s", err)
		}
	}

	for _, persona := range peeps {
		if exists, err := repo.Exists(context.Background(), db, persona.ID); !exists || err != nil {
			t.Errorf("unexpected Exist(): %v, err=%s", exists, err)
		}
	}
}

func TestPersonaRepository_Update_REAL(t *testing.T) {
	db := initTestDB(testDbDrv, testDbConn)
	defer db.Close()
	defer deinitSqlite3(testDbConn)

	testPersona := types.Persona{
		ID:   1,
		Name: "X-ray Bradbury",
		DOB:  mustPtrTime("1920-08-22 00:00:01"),
	}

	repo := NewPersonaRepository()
	check(repo.Insert(context.Background(), db, &testPersona))

	testPersona.Name = "Ray Bradbury"
	if err := repo.Update(context.Background(), db, &testPersona); err != nil {
		// Will get an failure on a real database
		t.Fatalf("Update() failed: %s", err)
	}
}
