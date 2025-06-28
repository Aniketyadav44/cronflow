package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func loadDb(pgHost, pgPort, pgUser, pgPass, pgDbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
