package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func ValidSession(t *testing.T, e *jtesting.E, c *Connection) *types.Session {
	session := new(types.Session)
	session.CreatedAt = timestamppb.Now()
	session.ExpiresAt = &timestamppb.Timestamp{
		Seconds: 1<<57 - 1,
	}
	t.Log(session)

	err := c.DoTx(e.Context, e.Logger, false, func(tx *sqlx.Tx) error {
		user := &types.User{
			Username: "test",
		}
		if err := UpdateUser(e.Context, tx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		session.User = user
		id, err := sessions.GenerateID()
		if err != nil {
			return fmt.Errorf("generate session id: %w", err)
		}
		session.Id = id
		if err := AddSession(e.Context, tx, session); err != nil {
			return fmt.Errorf("add session: %w", err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("create valid session: %v", err)
	}
	return session
}
