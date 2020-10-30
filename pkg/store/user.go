package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/types"
)

// LookupUser fills in the provided user object, searching by ID or Username.
func LookupUser(ctx context.Context, db sqlx.ExtContext, user *types.User) error {
	if id := user.GetId(); id != 0 {
		row := db.QueryRowxContext(ctx, `select username from "user" where id=$1`, id)
		if err := row.StructScan(user); err != nil {
			return fmt.Errorf("get user by id: %w", err)
		}
		return nil
	}
	if username := user.GetUsername(); username != "" {
		row := db.QueryRowxContext(ctx, `select id from "user" where username=$1`, username)
		if err := row.StructScan(user); err != nil {
			return fmt.Errorf("get user by username: %w", err)
		}
		return nil
	}
	return &ErrEmpty{Field: "(oneof:user.id,user.username)"}
}

// UpdateUser edits the provided user, creating it if it doesn't exist.
func UpdateUser(ctx context.Context, db sqlx.ExtContext, user *types.User) error {
	if user.Username == "" {
		return &ErrEmpty{Field: "username"}
	}
	if user.Id == 0 {
		rows, err := sqlx.NamedQueryContext(ctx, db, `insert into "user" (username) values (:username) returning (id)`, user)
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}
		defer rows.Close()
		if ok := rows.Next(); !ok {
			return errors.New("insert: no id returned")
		}
		if err := rows.Scan(&user.Id); err != nil {
			return fmt.Errorf("insert: scan id: %w", err)
		}
		return nil
	}

	info, err := sqlx.NamedExecContext(ctx, db, `update "user" set username=:username where id=:id`, user)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	affected, err := info.RowsAffected()
	if err != nil {
		return fmt.Errorf("update: get affected rows: %w", err)
	}
	if got, want := affected, int64(1); got != want {
		if got == 0 {
			return ErrNothingToUpdate
		}
		return fmt.Errorf("update: affected rows: got %v want %v", got, want)
	}
	return nil
}
