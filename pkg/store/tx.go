package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	txStarted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jsso2_db_tx_started",
		Help: "Number of attempts to begin a transaction via DoTx.",
	})
	txFinished = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "jsso2_db_tx_finished",
		Help: "Number of transactions started via DoTx that finished.",
	}, []string{"status"})
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
func (c *Connection) DoTx(origCtx context.Context, l *zap.Logger, readOnly bool, f func(tx *sqlx.Tx) error) error {
	span, ctx := opentracing.StartSpanFromContext(origCtx, "do_tx")
	var errs []error
	defer span.Finish()
	defer func() {
		if len(errs) > 0 {
			// We only log this to avoid suppressing errors that might be interesting to
			// a developer or operator.  You'll see this if a transaction succeeds after
			// a failed attempt, but sometimes the extra information is interesting.
			l.Debug("DoTx: early return hides some errors", zap.Errors("errors", errs))
		}
	}()
	for i := 0; i < MaxRetries; i++ {
		span.LogKV("attempt", i)
		if i != 0 {
			l.Debug("DoTx: retrying transaction after a delay", zap.Int("attempt", i), zap.Int("max_attempts", MaxRetries), zap.Errors("errors", errs), zap.Duration("delay", TxDelay))
			time.Sleep(TxDelay)
		}
		txStarted.Inc()
		tx, err := c.db.BeginTxx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  readOnly,
		})
		if err != nil {
			txFinished.WithLabelValues("failed_start").Inc()
			if isRetryable(err) {
				errs = append(errs, fmt.Errorf("attempt %d: begin tx: %w", i, err))
				continue
			}
			return fmt.Errorf("attempt %d: begin tx: non-retryable error: %w", i, err)
		}
		if err := f(tx); err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				if !errors.Is(rErr, sql.ErrTxDone) {
					txFinished.WithLabelValues("failed_rollback").Inc()
					return fmt.Errorf("attempt %d: error in user function (%v), but rollback failed too: %w", i, err, rErr)
				}
			}
			txFinished.WithLabelValues("rollback").Inc()
			if isRetryable(err) {
				errs = append(errs, fmt.Errorf("attempt %d: user function: %w", i, err))
				continue
			}
			return fmt.Errorf("attempt %d: user function: non-retryable error: %w", i, err)
		}
		if err := tx.Commit(); err != nil {
			txFinished.WithLabelValues("failed_commit").Inc()
			if isRetryable(err) {
				errs = append(errs, fmt.Errorf("attempt %d: commit: %w", i, err))
				continue
			}
		}
		txFinished.WithLabelValues("commit").Inc()
		return nil
	}
	ext.Error.Set(span, true)
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
