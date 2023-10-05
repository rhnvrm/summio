package database

import (
	"fmt"
	"time"
)

func (m *Manager) CheckMigrationTableExists() error {
	rows, err := m.db.Query("SELECT * FROM migration_meta")
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (m *Manager) CheckVersionApplied(version string) (bool, error) {
	rows, err := m.db.Query("SELECT * FROM migration_meta WHERE version = $1", version)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func (m *Manager) RunMigration(version, sql string) error {
	_, err := m.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("could not execute migration: %v", err)
	}

	// update the migration table
	_, err = m.db.Exec(
		"INSERT INTO migration_meta (version, applied_at) VALUES ($1, $2)",
		version,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("could not update migration table: %v", err)
	}

	return nil
}
