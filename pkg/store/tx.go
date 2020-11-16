package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Retryable allows you to explicitly mark an error as retryable.
type Retryable struct {
	Err error
}

func (r *Retryable) Error() string {
	return r.Err.Error()
}

func WrapRetryable(err error) error {
	return &Retryable{Err: err}
}

func (r *Retryable) Unwrap() error {
	return r.Err
}

// isRetryable determines whether or not an error can be fixed by retrying the transaction.
func isRetryable(err error) bool {
	target := &Retryable{}
	if errors.As(err, &target) {
		return true
	}
	if errors.Is(err, sql.ErrTxDone) {
		return true
	}
	if strings.Contains(err.Error(), "SQLSTATE 57") { // "terminating connection due to administrator command", etc.
		return true
	}
	return false
}

// DoTx executes the provied function in a transaction, retrying it if it rolls back.  You should
// not manually commit or roll back the provided transaction; return an error to roll back or return
// nil to commit.
func (c *Connection) DoTx(ctx context.Context, l *zap.Logger, readOnly bool, f func(tx *sqlx.Tx) error) error {
	var errs []error
	defer func() {
		if len(errs) > 0 {
			// We only log this to avoid suppressing errors that might be interesting to
			// a developer or operator.  You'll see this if a transaction succeeds after
			// a failed attempt, but sometimes the extra information is interesting.
			l.Debug("DoTx: early return hides some errors", zap.Errors("errors", errs))
		}
	}()
	for i := 0; i < MaxRetries; i++ {
		if i != 0 {
			l.Debug("DoTx: retrying transaction after a delay", zap.Int("attempt", i), zap.Int("max_attempts", MaxRetries), zap.Errors("errors", errs), zap.Duration("delay", TxDelay))
			time.Sleep(TxDelay)
		}
		tx, err := c.db.BeginTxx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  readOnly,
		})
		if err != nil {
			if isRetryable(err) {
				errs = append(errs, fmt.Errorf("attempt %d: begin tx: %w", i, err))
				continue
			}
			return fmt.Errorf("attempt %d: begin tx: non-retryable error: %w", i, err)
		}
		if err := f(tx); err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				if !errors.Is(rErr, sql.ErrTxDone) {
					return fmt.Errorf("attempt %d: error in user function (%v), but rollback failed too: %w", i, err, rErr)
				}
			}
			if isRetryable(err) {
				errs = append(errs, fmt.Errorf("attempt %d: user function: %w", i, err))
				continue
			}
			return fmt.Errorf("attempt %d: user function: non-retryable error: %w", i, err)
		}
		if err := tx.Commit(); err != nil {
			if isRetryable(err) {
				errs = append(errs, fmt.Errorf("attempt %d: commit: %w", i, err))
				continue
			}
		}
		return nil
	}
	if n := len(errs); n == 0 {
		l.Panic("DoTx: no errors, but loop ended; that's impossible")
		// Go doesn't know that this panics, so we pretend to return an error.
		return errors.New("bug: no errors?")
	} else if n == 1 {
		err := errs[0]
		errs = nil
		return err
	} else {
		// We do this dance so that Unwrap on the returned error yields the error from the last attempt.
		msg := new(strings.Builder)
		msg.WriteString(fmt.Sprintf("transaction failed after %d attempts: \n", MaxRetries))
		for i := 0; i < len(errs)-1; i++ {
			msg.WriteString("    ")
			msg.WriteString(errs[i].Error())
			msg.WriteString("\n")
		}
		err := errs[len(errs)-1]
		errs = nil
		return fmt.Errorf("%s    %w", msg.String(), err)
	}
}
