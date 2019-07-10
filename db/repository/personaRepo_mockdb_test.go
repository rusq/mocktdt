package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rusq/mocktdt/types"
)

func TestPersonaRepository_Update_MOCK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("something is terribly wrong: %s", err)
	}

	mock.ExpectExec(`INSERT`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE.*SET.*WHERE`).WillReturnResult(sqlmock.NewResult(1, 1))

	testPersona := types.Persona{
		ID:   1,
		Name: "X-ray Bradbury",
		DOB:  mustPtrTime("1920-08-22 00:00:01"),
	}
	// ...
	repo := NewPersonaRepository()
	// this is not required when using mock, leaving just to keep the same struct
	check(repo.Insert(context.Background(), db, &testPersona))
	testPersona.Name = "Ray Bradbury"
	if err := repo.Update(context.Background(), db, &testPersona); err != nil {
		// Will get an failure on a real database
		t.Errorf("Update() failed: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
