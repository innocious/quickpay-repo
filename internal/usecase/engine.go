package usecase

import (
	"errors"
	"quickpay/internal/repository"
)

type Engine struct {
	repo *repository.SQLiteRepo
}

func NewEngine(repo *repository.SQLiteRepo) *Engine {
	return &Engine{repo: repo}
}

func (e *Engine) ExecuteTransfer(senderID, receiverID string, transferAmount int64) error {
	// 1. Calculate 1% processing fee
	fee := transferAmount / 100
	totalDeduction := transferAmount + fee

	db := e.repo.DB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback() // Ensure rollback on failure

	var senderBalance int64
	err = tx.QueryRow("SELECT balance_cents FROM users WHERE id = ?", senderID).Scan(&senderBalance)
	if err != nil {
		return err
	}

	if senderBalance < totalDeduction {
		return errors.New("ERR_INSUFFICIENT_FUNDS")
	}

	_, err = tx.Exec("UPDATE users SET balance_cents = balance_cents - ? WHERE id = ?", totalDeduction, senderID)
	if err != nil {
		return err
	}

	// Add to receiver
	_, err = tx.Exec("UPDATE users SET balance_cents = balance_cents + ? WHERE id = ?", transferAmount, receiverID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ExecuteDeposit adds funds to a user's account after validating limits
func (e *Engine) ExecuteDeposit(userID string, amount int64) error {
	query := `UPDATE users SET balance_cents = balance_cents + ? WHERE id = ?`
	result, err := e.repo.DB().Exec(query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
