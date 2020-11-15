package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/opinionated-server/server"
	"go.uber.org/zap"
)

func main() {
	server.AppName = "jsso2-regenerate-db"
	dbConfig := &store.Config{}
	server.AddFlagGroup("database", dbConfig)
	server.Setup()

	ctx, c := context.WithTimeout(context.Background(), time.Minute)

	cfg, err := pgx.ParseConfig(dbConfig.DatabaseURL)
	if err != nil {
		zap.L().Fatal("problem parsing database url", zap.String("database_url", dbConfig.DatabaseURL), zap.Error(err))
	}

	name := cfg.Database
	if name == "" {
		zap.L().Fatal("no database to delete")
	}
	cfg.Database = "postgres"

	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		zap.L().Fatal("problem connecting to database", zap.String("connection_string", cfg.ConnString()), zap.Error(err))
	}

	zap.L().Info("deleting database", zap.String("database_name", name))
	if _, err := conn.Exec(ctx, fmt.Sprintf("drop database if exists %q with (force)", name)); err != nil {
		zap.L().Fatal("problem deleting database", zap.String("database_name", name), zap.Error(err))
	}

	zap.L().Info("creating database", zap.String("database_name", name))
	if _, err := conn.Exec(ctx, fmt.Sprintf("create database %q", name)); err != nil {
		zap.L().Fatal("problem creating new database", zap.String("database_name", name), zap.Error(err))
	}
	conn.Close(ctx)

	if dbConfig.RunMigrations {
		db, err := store.Connect(ctx, dbConfig.DatabaseURL)
		if err != nil {
			zap.L().Fatal("problem connecting to database", zap.String("database_url", dbConfig.DatabaseURL), zap.Error(err))
		}
		zap.L().Info("running database migrations")
		if err := db.MigrateDB(ctx); err != nil {
			zap.L().Fatal("problem running database migrations", zap.Error(err))
		}
	}
	c()
}
