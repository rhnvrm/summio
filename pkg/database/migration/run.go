package migration

import (
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/rhnvrm/summio/pkg/database"
)

// Embed all the migration SQL files in the app.
//
//go:embed *.sql
var migrationsFS embed.FS

func RunMigration(dbPath string) {
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("could not init db: %v", err)
	}

	m := NewMigrationManager(db)

	// List all the migration files.
	dirEntry, err := migrationsFS.ReadDir(".")
	if err != nil {
		log.Fatalf("could not read migrations directory: %v", err)
	}

	// sort dirEntry by name
	sort.Slice(dirEntry, func(i, j int) bool {
		return dirEntry[i].Name() < dirEntry[j].Name()
	})

	for _, entry := range dirEntry {
		fname := entry.Name()
		version := strings.Split(fname, "_")[0]
		// Read the file contents.
		fileContents, err := migrationsFS.ReadFile(fname)
		if err != nil {
			log.Fatalf("could not read migration file: %v", err)
		}

		// Execute the migration.
		if err := m.executeMigration(version, string(fileContents)); err != nil {
			log.Fatalf("could not execute migration: %v", err)
		}
	}
}

type MigrationManager struct {
	db *database.Manager
}

func NewMigrationManager(db *database.Manager) *MigrationManager {
	return &MigrationManager{db: db}
}

func (m *MigrationManager) executeMigration(version, sql string) error {
	var skipMigrationCheck bool

	// check if migration table exists
	if err := m.db.CheckMigrationTableExists(); err != nil {
		skipMigrationCheck = true
	}

	// check if migration was already applied
	if !skipMigrationCheck {
		applied, err := m.db.CheckVersionApplied(version)
		if err != nil {
			return err
		}

		if applied {
			log.Printf("migration %s already applied", version)
			return nil
		}
	}

	log.Printf("executing migration %s: %s", version, sql)
	if err := m.db.RunMigration(version, sql); err != nil {
		return fmt.Errorf("could not run migration: %v", err)
	}

	return nil
}
