package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"quickpay/internal/domain"
	"quickpay/internal/repository"
	"quickpay/internal/usecase"
)

// This is the entry point of the QuickPay API server.
// It initializes the database connection, performs migrations, sets up HTTP routes for user creation and money transfers, and starts the server on port 8080.
// The server handles JSON payloads for both endpoints and returns appropriate HTTP status codes based on the success or failure of the operations.
func main() {
	log.Println("INIT: Starting main function...")

	repo, err := repository.NewSQLiteRepository("quickpay.db")
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}

	defer repo.Close()

	log.Println("INIT: database connected successfully")

	if err := repo.Migrate(); err != nil {
		log.Fatalf("FATAL: Failed to migrate database: %v", err)
	}

	log.Println("INIT: database migrated successfully")

	engine := usecase.NewEngine(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		var u domain.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if err := domain.ValidateUser(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save to the SQLite DB
		if err := repo.CreateUser(u); err != nil {
			// Log the actual error to your internal terminal, but don't send it to the user
			log.Printf("CreateUser failed: %v", err)

			// Map the domain error to an HTTP response
			if errors.Is(err, domain.ErrDuplicateUser) {
				http.Error(w, `{"status": "error", "message": "Account already exists"}`, http.StatusConflict)
				return
			}

			// Catch-all for true internal database crashes
			http.Error(w, `{"status": "error", "message": "Internal Server Error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "success", "message": "User created successfully"}`))
	})

	mux.HandleFunc("POST /api/transfer", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			SenderID   string `json:"sender_id"`
			ReceiverID string `json:"receiver_id"`
			Amount     int64  `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		err := engine.ExecuteTransfer(payload.SenderID, payload.ReceiverID, payload.Amount)
		if err != nil {
			if err.Error() == "ERR_INSUFFICIENT_FUNDS" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "Transfer completed successfully"}`))
	})

	mux.HandleFunc("POST /api/deposit", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			UserID string `json:"user_id"`
			Amount int64  `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if err := engine.ExecuteDeposit(payload.UserID, payload.Amount); err != nil {
			log.Printf("Deposit failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success", "message": "Funds deposited successfully"}`))
	})

	log.Println("Server booting... QuickPay API live on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
