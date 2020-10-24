package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/jtesting"
)

func MustGetTestDB(t *testing.T, e *jtesting.E) *Connection {
	t.Helper()
	ctx := e.Context
	c, err := Wrap(ctx, e.DB)
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}
	if err := c.MigrateDB(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return c
}

var ErrUnimplemented = errors.New("unimplemented")

type EmptyDB struct{}

func (EmptyDB) DriverName() string     { return "" }
func (EmptyDB) Rebind(_ string) string { return "" }
func (EmptyDB) BindNamed(_ string, _ interface{}) (string, []interface{}, error) {
	return "", nil, ErrUnimplemented
}
func (EmptyDB) QueryContext(_ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return nil, ErrUnimplemented
}
func (EmptyDB) QueryxContext(_ context.Context, _ string, _ ...interface{}) (*sqlx.Rows, error) {
	return nil, ErrUnimplemented
}
func (EmptyDB) QueryRowxContext(_ context.Context, _ string, _ ...interface{}) *sqlx.Row {
	return &sqlx.Row{}
}

type EmptyResult struct{}

func (EmptyResult) LastInsertId() (int64, error) { return 0, ErrUnimplemented }
func (EmptyResult) RowsAffected() (int64, error) { return 0, ErrUnimplemented }
func (EmptyDB) ExecContext(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return EmptyResult{}, ErrUnimplemented
}
