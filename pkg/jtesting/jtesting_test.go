package jtesting

import (
	"database/sql"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogging(t *testing.T) {
	if got, want := zap.L().Core().Enabled(zapcore.DebugLevel), false; got != want {
		t.Fatalf("global logger in an unexpected state: debug level enabled:\n  got: %v\n want: %v", got, want)
	}
	Run(t, "logging", R{Logger: true}, func(t *testing.T, extras *Extras) {
		if extras.Logger == nil {
			t.Fatal("didn't get a logger")
		}
		extras.Logger.Info("this is the per-test logger")
		zap.L().Info("this is the global logger")
	})
	if got, want := zap.L().Core().Enabled(zapcore.DebugLevel), false; got != want {
		t.Fatalf("after test: debug level enabled:\n  got: %v\n want: %v", got, want)
	}
}

func TestDatabase(t *testing.T) {
	Run(t, "db", R{Database: true}, func(t *testing.T, extras *Extras) {
		if extras.DSN == "" {
			t.Fatal("DSN is empty")
		}
		c, err := sql.Open("pgx", extras.DSN)
		if err != nil {
			t.Fatalf("connect to db %q: %v", extras.DSN, err)
		}
		c.Close()
	})
}
