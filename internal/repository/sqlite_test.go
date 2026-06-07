package repository

import (
	"testing"
	"database/sql"
)

func TestDatabaseLiveness(t *testing.T) {
	// ARRANGE: Initialize our repository pointing to a temporary test database
	repo, err := NewSQLiteRepository("file:test_liveness_db?mode=memory")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close() // Clean up after the test

	// ACT: Ping the database to ensure it is actually live
	err = repo.Ping()

	// ASSERT
	if err != nil {
		t.Errorf("Expected database to be live, but ping failed: %v", err)
	}
}

func TestDatabaseMigration(t *testing.T) {
	// ARRANGE
	repo, err := NewSQLiteRepository("file:test_migration_db?mode=memory")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	
	defer repo.Close()

	// ACT
	err = repo.Migrate()
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// ASSERT
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='users';`
	err = repo.db.QueryRow(query).Scan(&tableName)
	
	if err == sql.ErrNoRows {
		t.Errorf("Expected 'users' table to exist, but it was not found.")
	} else if err != nil {
		t.Fatalf("Failed to query sqlite_master: %v", err)
	}
}


	