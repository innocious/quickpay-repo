package domain

import "testing"

func TestValidateDeposit_UnderLimit(t *testing.T) {
	// ARRANGE: $4.99 deposit (499 cents)
	tx := Transaction{
		UserID: "user_123",
		Amount: 499,
	}

	//ACT
	err := ValidateDeposit(tx)

	//ASSERT
	if err == nil || err.Error() != "ERR_DEPOSIT_LIMIT" {
		t.Errorf("Expected ERR_DEPOSIT_LIMIT for $4.99, got: %v", err)
	}
}

func TestValidateDeposit_OverLimit(t *testing.T) {
	// ARRANGE: $5,000.01 deposit (500001 cents)
	tx := Transaction{
		UserID: "user_123",
		Amount: 500001,
	}

	// ACT
	err := ValidateDeposit(tx)

	// ASSERT
	if err == nil || err.Error() != "ERR_DEPOSIT_LIMIT" {
		t.Errorf("Expected ERR_DEPOSIT_LIMIT for $5,000.01, got: %v", err)
	}
}
