package domain

import (
	"testing"
)

func TestValidateUser_AgeRestricted(t *testing.T) {
	// ARRANGE
	newUser := User{
		ID:        "user_123",
		LegalName: "Test Account",
		Email:     "test@example.com",
		Age:       17,
	}

	// ACT
	err := ValidateUser(newUser)

	// ASSERT
	if err == nil {
		t.Fatalf("Expected an error for an underage user, but got nil")
	}

	if err.Error() != "ERR_AGE_RESTRICTED" {
		t.Errorf("Expected error message 'ERR_AGE_RESTRICTED', got '%s'", err.Error())
	}
}

func TestValidateUser_ValidateEmail(t *testing.T) {
	// ARRANGE
	newUser := User{
		ID:        "user_124",
		LegalName: "Test Account",
		Email:     "ivalidemail.com",
		Age:       25,
	}

	// ACT
	err := ValidateUser(newUser)

	// ASSERT
	if err == nil {
		t.Fatalf("Expected an error for an invalid email, but got nil")
	}

	if err.Error() != "ERR_INVALID_EMAIL" {
		t.Errorf("Expected error message 'ERR_INVALID_EMAIL', got '%s'", err.Error())
	}

}
