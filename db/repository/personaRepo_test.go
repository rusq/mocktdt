package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rusq/mocktdt/types"
)

func TestPersonaRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("something is terribly wrong: %s", err)
	}

	type args struct {
		ctx     context.Context
		db      Execer
		persona *types.Persona
	}
	tests := []struct {
		name         string
		expectations func(sqlmock.Sqlmock)
		args         args
		wantErr      bool
	}{
		{"OK",
			func(m sqlmock.Sqlmock) {
				m.ExpectExec(`INSERT.*VALUES.*`).
					WithArgs(1, "William Blazkowicz", nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args{context.Background(), db, &types.Persona{ID: 1, Name: "William Blazkowicz"}},
			false},
		{"error",
			func(m sqlmock.Sqlmock) {
				m.ExpectExec(`INSERT.*VALUES.*`).
					WithArgs(2, "Stan Blazkowicz", nil).
					WillReturnError(errors.New("too many Blazkowicz in the DB"))
			},
			args{context.Background(), db, &types.Persona{ID: 2, Name: "Stan Blazkowicz"}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PersonaRepository{}
			tt.expectations(mock)
			if err := p.Insert(tt.args.ctx, tt.args.db, tt.args.persona); (err != nil) != tt.wantErr {
				t.Errorf("PersonaRepository.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
