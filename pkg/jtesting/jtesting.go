// Package jtesting contains test helpers for JSSO.  (It's called jtesting so you don't have to alias an import.)
package jtesting

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // This is the only driver we support.
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Config contains test-specific configuration.
type Config struct {
	SuperuserDSN string // The DSN to use to create databases.
}

func makeConfig() *Config {
	return &Config{
		SuperuserDSN: "user=postgres host=localhost port=5432 sslmode=disable",
	}
}

// R requests specific extras during a test run.
type R struct {
	Timeout  time.Duration
	Logger   bool
	Database bool
}

// Extras holds per-test "extras".
type Extras struct {
	Context context.Context
	Logger  *zap.Logger
	Config  *Config
	DSN     string
	DB      *sql.DB
}

// Run runs the provided test function as a subtest with the desired Extras available.
func Run(t *testing.T, name string, r R, f func(t *testing.T, extras *Extras)) {
	t.Helper()
	pc, _, _, pcOk := runtime.Caller(1)
	t.Run(name, func(t *testing.T) {
		var ctx context.Context
		var c func()
		extras := &Extras{
			Config: makeConfig(),
		}
		if r.Timeout > 0 {
			ctx, c = context.WithTimeout(context.Background(), r.Timeout)
		} else {
			ctx, c = context.WithCancel(context.Background())
		}
		extras.Context = ctx
		if r.Logger {
			logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
			defer logger.Sync()
			restoreLogger := zap.ReplaceGlobals(logger.Named("global"))
			defer restoreLogger()
			extras.Logger = logger.Named("test." + name)
		}
		if r.Database {
			if !pcOk {
				t.Fatal("could not determine caller to generate database name")
			}
			dsn, err := newTestDB(ctx, pc, name, extras.Config.SuperuserDSN)
			if err != nil {
				t.Fatalf("creating test database: %v", err)
			}
			extras.DSN = dsn
			db, err := sql.Open("pgx", dsn)
			if err != nil {
				t.Fatalf("connect to test database: %v", err)
			}
			extras.DB = db
		}
		f(t, extras)
		select {
		case <-ctx.Done():
			t.Fatalf("after tests: %v", ctx.Err())
		default:
		}
		c()
	})
}

// newTestDB creates a new test database.
func newTestDB(ctx context.Context, pc uintptr, name, databaseURL string) (string, error) {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "", fmt.Errorf("cannot determine database name from caller: pc %v does not map to a function", pc)
	}
	name = fmt.Sprintf("jsso-test-%s-%s", f.Name(), name)
	cfg, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return "", fmt.Errorf("parse databse url: %w", err)
	}
	c, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		return "", fmt.Errorf("connect %s: %w", cfg.ConnString(), err)
	}
	defer c.Close()
	if _, err := c.ExecContext(ctx, fmt.Sprintf("drop database if exists %q", name)); err != nil {
		return "", fmt.Errorf("drop old database %s: %w", name, err)
	}
	if _, err := c.ExecContext(ctx, fmt.Sprintf("create database %q", name)); err != nil {
		return "", fmt.Errorf("create database %s: %w", name, err)
	}
	cfg.Database = name
	return cfg.ConnString(), nil
}
