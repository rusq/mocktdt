package repository

import (
	"context"
	"database/sql"

	"github.com/rusq/mocktdt/types"
)

// PersonaRepository is a repository for Persona type.
type PersonaRepository struct{}

// NewPersonaRepository creates a new persona repository.
func NewPersonaRepository() *PersonaRepository {
	return &PersonaRepository{}
}

// Execer is an interface for Exec* functions.
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Getter is an interface for Query* functions.
type Getter interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Insert inserts Persona into the database.
func (*PersonaRepository) Insert(ctx context.Context, db Execer, persona *types.Persona) error {
	stmt := `INSERT INTO PERSONAS (ID, NAME, DOB) VALUES (?, ?, ?)`
	_, err := db.ExecContext(ctx, stmt, persona.ID, persona.Name, persona.DOB)
	return err
}

// Update updates the persona in the database.
func (*PersonaRepository) Update(ctx context.Context, db Execer, persona *types.Persona) error {
	// ORA-00917: missing comma
	stmt := `UPDATE PERSONAS SET NAME=? DOB=? WHERE ID=?`
	_, err := db.ExecContext(ctx, stmt, persona.Name, persona.DOB, persona.ID)
	return err
}

// Delete deletes persona from the database.
func (*PersonaRepository) Delete(ctx context.Context, db Execer, id int) error {
	stmt := `DELETE FROM PERSONAS WHERE ID=?`

	_, err := db.ExecContext(ctx, stmt, id)
	return err
}

// Exists checks if the persona exists in the database.
func (*PersonaRepository) Exists(ctx context.Context, db Getter, id int) (bool, error) {
	stmt := "SELECT COALESCE(MAX(ID),0) FROM PERSONAS WHERE ID=?"

	row := db.QueryRowContext(ctx, stmt, id)
	var exists int64
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists > 0, nil
}
