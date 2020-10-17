package store

import (
	"testing"

	"github.com/jrockway/jsso2/pkg/jtesting"
)

func MustGetTestDB(t *testing.T, e *jtesting.E) *Connection {
	t.Helper()
	ctx := e.Context
	c, err := Wrap(ctx, e.DB)
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}
	if err := c.MigrateDB(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return c
}
