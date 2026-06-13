# QuickPay Core Engine

A lightweight, strictly validated financial transaction engine built in Go. This system implements automated data validation, boundary constraints, and ACID-compliant atomic database transfers using an embedded SQLite ledger.

## Prerequisites
- Go 1.22 or higher
- `curl` (or Postman) for triggering API requests

## Setup & Execution

This application requires no external database configuration. It automatically provisions a local SQLite database (`quickpay.db`) and runs migrations on boot.

1. **Start the Server:**
   Navigate to the root directory of the project and execute:
   ```bash
   go run ./cmd/api/main.go