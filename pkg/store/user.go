package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/types"
)

// UpdateUser edits the provided user, creating it if it doesn't exist.
func (c *Connection) UpdateUser(ctx context.Context, tx *sqlx.Tx, user *types.User) error {
	if user.Id == 0 {
		rows, err := sqlx.NamedQueryContext(ctx, tx, `insert into "user" (username) values (:username) returning (id)`, user)
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

	info, err := sqlx.NamedExecContext(ctx, tx, `update "user" set username=:username where id=:id`, user)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	affected, err := info.RowsAffected()
	if err != nil {
		return fmt.Errorf("update: get affected rows: %w", err)
	}
	if got, want := affected, int64(1); got != want {
		return fmt.Errorf("update: affected rows: got %v want %v", got, want)
	}
	return nil
}
