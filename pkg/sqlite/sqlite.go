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

	data, _ := os.ReadFile("./database/schemes/up_tables.sql")

	if _, err := db.Exec(string(data)); err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	// if err := createTables(db, cfg.Sqlite.SchemePath); err != nil {
	// 	return nil, fmt.Errorf("connect db: %w", err)
	// }

	return db, nil
}

func createTables(db *sql.DB, schemePath string) error {
	schemes, err := readSchemes(schemePath)
	if err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	for _, scheme := range schemes {
		stmt, err := db.Prepare(scheme)
		if err != nil {
			return fmt.Errorf("create tables: prepare: %w", err)
		}

		if _, err := stmt.Exec(); err != nil {
			return fmt.Errorf("create tables: exec: %w", err)
		}

		stmt.Close()
	}

	return nil
}

func readSchemes(schemePath string) ([]string, error) {
	var schemes []string

	files, err := os.ReadDir(schemePath)
	if err != nil {
		return nil, fmt.Errorf("read schemes: read dir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			temp, err := os.ReadFile(schemePath + file.Name())
			if err != nil {
				return nil, fmt.Errorf("read schemes: read file: %w", err)
			}

			schemes = append(schemes, string(temp))
		}
	}

	return schemes, nil
}
