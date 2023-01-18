package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"real-time-forum/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase(cfg *config.Config) (*sql.DB, error) {
	enableForeignKeys := "?_foreign_keys=on&cache=shared&mode=rwc"

	db, err := sql.Open(cfg.Sqlite.Driver, cfg.Sqlite.DatabaseFileName+enableForeignKeys)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := createTables(db, cfg.Sqlite.SchemePath); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB, schemePath string) error {
	data, err := os.ReadFile(schemePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	stmt, err := db.Prepare(string(data))
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}

	if _, err := stmt.Exec(); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	defer stmt.Close()

	return nil
}
