package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // This is the only driver we support.
	"github.com/jmoiron/sqlx"
)

const (
	MaxRetries = 3
	TxDelay    = 10 * time.Millisecond
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
