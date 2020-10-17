package store

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/types"
)

func TestUpdateUser(t *testing.T) {
	jtesting.Run(t, "updateusers", jtesting.R{Logger: true, Database: true}, func(t *testing.T, e *jtesting.E) {
		c := MustGetTestDB(t, e)
		c.DoTx(e.Context, e.Logger, false, func(tx *sqlx.Tx) error {
			user := &types.User{Username: "foo"}
			if err := c.UpdateUser(e.Context, tx, user); err != nil {
				t.Fatalf("create user: %v", err)
			}
			if user.Id == 0 {
				t.Errorf("user id not updated in place: %#v", user)
			}
			return nil
		})
	})
}
