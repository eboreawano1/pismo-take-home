package repository

import (
	"fmt"
	"context"
	"database/sql"
	"pismo-take-home/internal/event"
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

func (repository *Repository) Save(ctx context.Context, event event.Event) error {
	_, error := repository.database.ExecContext(ctx, `
		INSERT INTO processed_events (
			event_id,
			event_type,
			tenant_id,
			status
			producer,
			payload,
			event_time,
			schema_version,
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8)
	`,
		event.EventId,
		event.EventType,
		event.TenantId,
		"READY_TO_DELIVER",
		event.Producer,
		event.Payload,
		event.EventTime,
		event.SchemaVersion,
	
	)

	if error != nil {
		return fmt.Errorf("Error while saving event: %w", error)
	}
		
	return nil
}

func (repository *Repository) Close() error {
	return repository.database.Close()
}

type Repository struct {
	database *sql.DB
}