// Package jtesting contains test helpers for JSSO.  (It's called jtesting so you don't have to alias an import.)
package jtesting

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // This is the only driver we support.
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Config contains test-specific configuration.
type Config struct {
	SuperuserDSN string // The DSN to use to create databases.
}

// R requests specific extras during a test run.
type R struct {
	Timeout  time.Duration
	Logger   bool
	Database bool
}

// E holds per-test "extras".
type E struct {
	Context context.Context
	Logger  *zap.Logger
	Config  *Config
	DSN     string
	DB      *sql.DB
}

// Run runs the provided test function as a subtest with the desired Extras available.
func Run(t *testing.T, name string, r R, f func(t *testing.T, e *E)) {
	t.Helper()
	pc, file, _, pcOk := runtime.Caller(1)
	t.Run(name, func(t *testing.T) {
		var ctx context.Context
		var c func()

		envFile := filepath.Clean(filepath.Join(file, "..", "..", "..", "env.test"))
		if err := godotenv.Load(envFile); err != nil {
			t.Fatalf("failed to load %s: %v", envFile, err)
		}

		extras := &E{
			Config: &Config{
				SuperuserDSN: os.Getenv("SUPERUSER_DATABASE_URL"),
			},
		}
		if r.Timeout > 0 {
			ctx, c = context.WithTimeout(context.Background(), r.Timeout)
		} else {
			ctx, c = context.WithCancel(context.Background())
		}
		if r.Logger {
			logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
			defer logger.Sync()
			restoreLogger := zap.ReplaceGlobals(logger.Named("global"))
			defer restoreLogger()
			extras.Logger = logger.Named("test." + name)
			ctx = ctxzap.ToContext(ctx, logger)
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
		extras.Context = ctx
		f(t, extras)
		select {
		case <-ctx.Done():
			t.Fatalf("after tests: %v", ctx.Err())
		default:
		}
		c()
	})
}

const pkgPrefix = "github.com/jrockway/jsso2/pkg/"

// newTestDB creates a new test database.
func newTestDB(ctx context.Context, pc uintptr, name, databaseURL string) (string, error) {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "", fmt.Errorf("cannot determine database name from caller: pc %v does not map to a function", pc)
	}

	// Name the database for the test.  Try very hard to keep it under 64 characters, so that
	// database names don't collide.
	candidate := fmt.Sprintf("%s-%s", f.Name(), name)
	if strings.HasPrefix(candidate, pkgPrefix) {
		candidate = candidate[len(pkgPrefix):]
	}
	name = fmt.Sprintf("jsso-test-%s", candidate)
	if len(name) > 64 {
		hash := md5.Sum([]byte(candidate))
		name = fmt.Sprintf("jsso-test-%x", hash)
	}
	name = strings.NewReplacer(`"`, ``, `'`, ``, ` `, `-`, `_`, `-`, `=`, `-`).Replace(name)

	cfg, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return "", fmt.Errorf("parse databse url: %w", err)
	}
	c, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		return "", fmt.Errorf("connect %s: %w", cfg.ConnString(), err)
	}
	defer c.Close()
	if _, err := c.ExecContext(ctx, fmt.Sprintf("drop database if exists %q with (force)", name)); err != nil {
		return "", fmt.Errorf("drop old database %s: %w", name, err)
	}
	if _, err := c.ExecContext(ctx, fmt.Sprintf("create database %q", name)); err != nil {
		return "", fmt.Errorf("create database %s: %w", name, err)
	}
	dsn := cfg.ConnString() + ` database=` + name
	cfg, err = pgx.ParseConfig(dsn)
	if err != nil {
		return "", fmt.Errorf("newly-created connect string is invalid: %v", err)
	}
	if got, want := cfg.Database, name; got != want {
		return "", fmt.Errorf("parsed database string is invalid: got %v want %v", got, want)
	}
	return dsn, nil
}
