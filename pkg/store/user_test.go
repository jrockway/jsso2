package store

import (
	"testing"

	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/types"
)

func TestUpdateUser(t *testing.T) {
	jtesting.Run(t, "updateusers", jtesting.R{Logger: true, Database: true}, func(t *testing.T, e *jtesting.E) {
		c := MustGetTestDB(t, e)

		testData := []struct {
			name                    string
			beforeCount, afterCount int
			f                       func(t *testing.T)
		}{
			{
				name:        "add user",
				beforeCount: 0,
				afterCount:  1,
				f: func(t *testing.T) {
					user := &types.User{Username: "foo"}
					if err := c.UpdateUser(e.Context, c.db, user); err != nil {
						t.Fatalf("create user: %v", err)
					}
					if user.Id == 0 {
						t.Errorf("user id not updated in place: %#v", user)
					}
				},
			},
			{
				name:        "add duplicate user",
				beforeCount: 1,
				afterCount:  1,
				f: func(t *testing.T) {
					user := &types.User{Username: "foo"}
					if err := c.UpdateUser(e.Context, c.db, user); err == nil {
						t.Errorf("expected error when adding a duplicate user")
					}
				},
			},
		}
		for _, test := range testData {
			t.Run(test.name, func(t *testing.T) {
				var n int
				if err := c.db.QueryRowxContext(e.Context, `select count(1) from "user"`).Scan(&n); err != nil {
					t.Fatalf("count: %v", err)
				}
				if got, want := n, test.beforeCount; got != want {
					t.Errorf("initial count:\n  got: %v\n want: %v", got, want)
				}

				test.f(t)

				if err := c.db.QueryRowxContext(e.Context, `select count(1) from "user"`).Scan(&n); err != nil {
					t.Fatalf("count: %v", err)
				}
				if got, want := n, test.afterCount; got != want {
					t.Errorf("count:\n  got: %v\n want: %v", got, want)
				}
			})
		}
	})
}
