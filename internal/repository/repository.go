package repository

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)


func New(databaseURL string) (*Repository, error) {
	database, error := sql.Open("postgres", databaseURL)

	if error != nil {
		return nil, fmt.Errorf("error opening database: %w", error)
	}

	pingError := database.Ping()

	if pingError != nil {
		return nil, fmt.Errorf("error pinging database: %w", pingError)
	}

	return &Repository{ database: database }, nil
}

func (repository *Repository) Close() error {
	return repository.database.Close()
}

type Repository struct {
	database *sql.DB
}