package repository

import (
	"testing"
	"database/sql"
	"quickpay/internal/domain"
)

func TestDatabaseLiveness(t *testing.T) {
	// ARRANGE: Initialize our repository pointing to a temporary test database
	repo, err := NewSQLiteRepository("file:test_liveness.db?mode=memory")
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
	repo, err := NewSQLiteRepository("file:test_migration.db?mode=memory")
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

func TestCreateUser(t *testing.T) {
	// ARRANGE
	repo, _ := NewSQLiteRepository("file:test_create_user.db?mode=memory")
	defer repo.Close()

	_ = repo.Migrate()

	newUser := domain.User{
		ID: "user_999",
		LegalName: "Maxwell Smart",
		Email: "maxwell@example.com",
		Age: 28,
	}

	// ACT
	err := repo.CreateUser(newUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// ASSERT
	var savedName string
	err = repo.db.QueryRow("SELECT legal_name FROM users WHERE id = ?", newUser.ID).Scan(&savedName)

	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}
	if savedName != "Maxwell Smart" {
		t.Errorf("Expected 'Maxwell Smart', got '%s'", savedName)
	}
}
