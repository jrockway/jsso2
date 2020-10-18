package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // This is the only driver we support.
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MaxRetries = 3
	TxDelay    = 10 * time.Millisecond
)

var (
	ErrNothingToUpdate = errors.New("nothing to update")
)

// Config is environment/command-line config for storage.
type Config struct {
	DatabaseURL   string `long:"database_url" description:"Postgres connection string pointing at the database" env:"DATABASE_URL"`
	RunMigrations bool   `long:"run_migrations" description:"If true, migrate the database after connecting." env:"RUN_MIGRATIONS"`
}

// Connection is a connection to storage for jsso.
type Connection struct {
	db *sqlx.DB
}

// Wrap wraps an existing connection to the database.
func Wrap(ctx context.Context, db *sql.DB) (*Connection, error) {
	c := &Connection{db: sqlx.NewDb(db, "pgx")}
	if err := c.db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("sqlx ping: %w", err)
	}
	return c, nil
}

// Connect connects to the database.
func Connect(ctx context.Context, dsn string) (*Connection, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx connect: %w", err)
	}
	return &Connection{db: db}, nil
}

// AsGRPCError converts a store error to one with a gRPC status code.  Is is valid to call with a
// nil error.
func AsGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNothingToUpdate) {
		return status.Error(codes.NotFound, err.Error())
	}
	if isRetryable(err) {
		// From codes: "Use Unavailable if the client can retry just the failing call."
		return status.Error(codes.Unavailable, err.Error())
	}
	// SQLSTATE 23XXX is a referential integrity violation; duplicate unique index, null where
	// the schema dictates non-null, etc.
	if strings.Contains(err.Error(), "(SQLSTATE 23") {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	return status.Error(codes.Unknown, err.Error())
}
