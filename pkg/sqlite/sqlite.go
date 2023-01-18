package sqlite

import (
	"database/sql"
	"fmt"
)

func ConnectDB(driver, dbName, schemePath string) (*sql.DB, error) {
	enableForeignKeys := "?_foreign_keys=on&cache=shared&mode=rwc"

	db, err := sql.Open(driver, dbName+enableForeignKeys)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := createTables(db, schemePath); err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB, schemePath string) error {
	schemes, err := readTables(schemePath)
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
