package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/tern/migrate"
)

const (
	versionTable  = "public.schema_version"
	migrationPath = "migrations"
)

func findMigrations() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	start, err := filepath.Abs(cwd)
	if err != nil {
		return "", fmt.Errorf("abs: %w", err)
	}
	parts := strings.Split(start, string(filepath.Separator))
	if len(parts) > 0 {
		// filepath.Join really wants a / there, not ""
		parts[0] = string(filepath.Separator)
	}

	var lookedAt []string
	for i := len(parts); i > 0; i-- {
		target := filepath.Join(filepath.Join(parts[0:i]...), migrationPath)
		lookedAt = append(lookedAt, target)
		info, err := os.Stat(target)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", err
		}
		if info.IsDir() {
			return target, nil
		}
	}
	return "", fmt.Errorf("no migrations found; looked in %v", lookedAt)
}

func (c *Connection) MigrateDB(ctx context.Context) error {
	conn, err := c.db.DB.Conn(ctx)
	if err != nil {
		return fmt.Errorf("get raw connection: %w", err)
	}
	err = conn.Raw(func(driverConn interface{}) error {
		conn := driverConn.(*stdlib.Conn).Conn()
		m, err := migrate.NewMigrator(ctx, conn, versionTable)
		if err != nil {
			return fmt.Errorf("new migrator: %w", err)
		}
		path, err := findMigrations()
		if err != nil {
			return fmt.Errorf("find migrations: %w", err)
		}
		if err := m.LoadMigrations(path); err != nil {
			return fmt.Errorf("load migrations from ./migrations: %w", err)
		}
		if err := m.Migrate(ctx); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
