package store

import (
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jrockway/jsso2/pkg/jtesting"
)

func TestDoTx(t *testing.T) {
	jtesting.Run(t, "updateusers", jtesting.R{Logger: true, Database: true}, func(t *testing.T, extras *jtesting.Extras) {
		c := MustGetTestDB(t, extras)

		// Don't retry non-retryable errors.
		var n int
		err := c.DoTx(extras.Context, extras.Logger, false, func(tx *sqlx.Tx) error {
			n++
			return errors.New("oh no")
		})
		if err == nil {
			t.Error("DoTx should have errored")
		}
		if got, want := n, 1; got != want {
			t.Errorf("retry count:\n  got: %v\n want: %v", got, want)
		}

		// Retry retryable errors.
		n = 0
		err = c.DoTx(extras.Context, extras.Logger, false, func(tx *sqlx.Tx) error {
			n++
			return WrapRetryable(errors.New("oh no"))
		})
		if err == nil {
			t.Error("DoTx should have errored")
		}
		if got, want := n, MaxRetries; got != want {
			t.Errorf("retry count:\n  got: %v\n want: %v", got, want)
		}

		// Retry if the transaction gets into the "already committed or rolled back" state.
		n = 0
		err = c.DoTx(extras.Context, extras.Logger, false, func(tx *sqlx.Tx) error {
			n++
			if n < 2 {
				tx.Rollback()
			}
			return nil
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got, want := n, 2; got != want {
			t.Errorf("retry count:\n  got: %v\n want: %v", got, want)
		}
	})
}
