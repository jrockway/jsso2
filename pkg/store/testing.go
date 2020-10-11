package store

import (
	"testing"

	"github.com/jrockway/jsso2/pkg/jtesting"
)

func MustGetTestDB(t *testing.T, extras *jtesting.Extras) *Connection {
	t.Helper()
	ctx := extras.Context
	c, err := Wrap(ctx, extras.DB)
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}
	if err := c.MigrateDB(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return c
}
