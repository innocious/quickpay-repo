package main

import (
	"encoding/json"
	"log"
	"net/http"

	"quickpay/internal/domain"
	"quickpay/internal/repository"
	"quickpay/internal/usecase"

)

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

		if err := repo.CreateUser(u); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "success", "message": "User created successfully"}`))
	})

	mux.HandleFunc("POST /api/transfer", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			SenderID     string `json:"sender_id"`
			ReceiverID   string `json:"receiver_id"`
			Amount int64  `json:"amount"`
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

	log.Println("Server booting... QuickPay API live on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}