package usecase

import (
	"testing"
	"quickpay/internal/domain"
	"quickpay/internal/repository"
)

func setupTestDB(t *testing.T) *repository.SQLiteRepo {
	repo, err := repository.NewSQLiteRepository("file::memory?cache=shared")
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	_ = repo.Migrate()

	return repo
}

func TestExecuteTransfer_HappyPath(t *testing.T) {
	// ARRANGE
	repo := setupTestDB(t)
	defer repo.Close()

	// Create two users
	_ = repo.CreateUser(domain.User{ID: "user_A", LegalName: "Alice", Email: "alice@test.com", Age: 25})
	_, _ = repo.DB().Exec(`UPDATE users SET balance_cents = 20000 WHERE id = 'user_A'`) // Give Alice $100

	_ = repo.CreateUser(domain.User{ID: "user_B", LegalName: "Bob", Email: "bob@test.com", Age: 30})

	engine := NewEngine(repo)
	
	// ACT
	err := engine.ExecuteTransfer("user_A", "user_B", 10000) // Transfer $100

	// ASSERT
	if err != nil {
		t.Fatalf("Expected transfer to succeed, but got error: %v", err)
	}

	var balanceA int64
	_ = repo.DB().QueryRow(`SELECT balance_cents FROM users WHERE id = 'user_A'`).Scan(&balanceA)
	if balanceA != 9900 {
		t.Errorf("Expected User A balance to be 9900 cents, got %d", balanceA)
	}

	var balanceB int64
	_ = repo.DB().QueryRow(`SELECT balance_cents FROM users WHERE id = 'user_B'`).Scan(&balanceB)
	if balanceB != 10000 {
		t.Errorf("Expected User B balance to be 10000 cents, got %d", balanceB)
	}
}

func TestExecuteTransfer_InsufficientFunds(t *testing.T) {
	// ARRANGE
	repo := setupTestDB(t)
	defer repo.Close()

	// Create two users
	_ = repo.CreateUser(domain.User{ID: "user_C", LegalName: "Charlie", Email: "charlie@test.com", Age: 25})
	_, _ = repo.DB().Exec(`UPDATE users SET balance_cents = 10000 WHERE id = 'user_C'`) // Give Charlie $100

	_ = repo.CreateUser(domain.User{ID: "user_D", LegalName: "Diana", Email: "diana@test.com", Age: 30})

	engine := NewEngine(repo)

	// ACT
	err := engine.ExecuteTransfer("user_C", "user_D", 10000) // Attempt to transfer $100

	// ASSERT
	if err == nil {
		t.Fatalf("Expected transfer to fail due to insufficient funds, but it succeeded")
	}
	if err.Error() != "ERR_INSUFFICIENT_FUNDS" {
		t.Errorf("Expected ERR_INSUFFICIENT_FUNDS, got: %v", err.Error())
	}

	var balanceC int64
	_ = repo.DB().QueryRow(`SELECT balance_cents FROM users WHERE id = 'user_C'`).Scan(&balanceC)
	if balanceC != 10000 {
		t.Errorf("Expected User C balance to remain 10000 cents, got %d", balanceC)
	}
}