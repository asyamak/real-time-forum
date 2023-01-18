package database

import (
	"database/sql"
	"real-time-forum/internal/config"
	"real-time-forum/pkg/database/sqlite"
)

type DataBaseBuilder interface {
	ConnectDatabase(cfg *config.Config) (*sql.DB, error)
}

func New(databaseName string) DataBaseBuilder {
	switch databaseName {
	case "sqlite":
		return sqlite.NewSqlite()
	case "postgres":
		return nil
	default:
		return nil
	}
}
