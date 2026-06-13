package usecase

import (
	"quickpay/internal/domain"
	"quickpay/internal/repository"
	"testing"
)

func setupTestDB(t *testing.T) *repository.SQLiteRepo {
	repo, err := repository.NewSQLiteRepository(":memory:")
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

func TestExecuteDeposit_HappyPath(t *testing.T) {
	// ARRANGE: Setup the strictly ephemeral memory database
	repo := setupTestDB(t)
	defer repo.Close()

	// Create our user starting with 0 balance
	_ = repo.CreateUser(domain.User{ID: "user_deposit_1", LegalName: "Alice", Email: "alice@test.com", Age: 25})

	engine := NewEngine(repo)

	// ACT: Deposit $500.00 (50000 cents)
	err := engine.ExecuteDeposit("user_deposit_1", 50000)

	// ASSERT: The transaction should succeed
	if err != nil {
		t.Fatalf("Expected deposit to succeed, got error: %v", err)
	}

	// Verify the ledger actually reflects the 50000 cents
	var finalBalance int64
	_ = repo.DB().QueryRow("SELECT balance_cents FROM users WHERE id = 'user_deposit_1'").Scan(&finalBalance)

	if finalBalance != 50000 {
		t.Errorf("Expected balance to be 50000, got %d", finalBalance)
	}
}

func TestExecuteDeposit_UserNotFound(t *testing.T) {
	// ARRANGE: Setup the strictly ephemeral memory database
	repo := setupTestDB(t)
	defer repo.Close()

	engine := NewEngine(repo)

	// ACT: Attempt to deposit into a non-existent account
	err := engine.ExecuteDeposit("ghost_user_999", 10000)

	// ASSERT: The engine must catch that the row wasn't updated
	if err == nil {
		t.Fatalf("Expected deposit to fail for non-existent user, but it succeeded")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got: %v", err)
	}
}
