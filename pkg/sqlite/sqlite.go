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
