package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

type Manager struct {
	db *sqlx.DB
}

func InitDB(path string) (*Manager, error) {
	db, err := sqlx.Connect("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	return &Manager{db: db}, nil
}
