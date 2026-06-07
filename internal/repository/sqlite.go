package repository

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type SQLiteRepo struct{
	db *sql.DB
}

func NewSQLiteRepository(dataSourceName string) (*SQLiteRepo, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteRepo{db: db}

	if err := repo.Ping(); err != nil {
		return nil, err
	}

	return repo, nil

}

func (r *SQLiteRepo) Ping() error {
	return r.db.Ping()
}

func (r *SQLiteRepo) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepo) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		legal_name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER NOT NULL,
		balance_cents INTEGER NOT NULL DEFAULT 0
	);`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}