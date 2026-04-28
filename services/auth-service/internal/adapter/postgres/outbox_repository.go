package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type OutboxRepository struct {
	db *pgxpool.Pool
}

func NewOutboxRepository(db *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (or *OutboxRepository) CreateEvent(ctx context.Context, event *domain.OutboxEvent) error {
	query := `INSERT INTO outbox (id, event_type, payload, processed, created_at) VALUES ($1, $2, $3, $4, $5)`

	if _, err := or.db.Exec(ctx, query, event.ID, event.EventType, event.Payload, event.Processed, event.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (or *OutboxRepository) GetUnprocessedEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	query := `SELECT id, event_type, payload, processed, created_at, processed_at 
			  FROM outbox 
			  WHERE processed = FALSE 
			  ORDER BY created_at ASC 
			  LIMIT $1`

	rows, err := or.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.OutboxEvent
	for rows.Next() {
		var event domain.OutboxEvent
		var processedAt *time.Time
		if err := rows.Scan(&event.ID, &event.EventType, &event.Payload, &event.Processed, &event.CreatedAt, &processedAt); err != nil {
			return nil, err
		}
		event.ProcessedAt = processedAt
		events = append(events, &event)
	}

	return events, nil
}

func (or *OutboxRepository) MarkAsProcessed(ctx context.Context, id string) error {
	now := time.Now().UTC()
	query := `UPDATE outbox SET processed = TRUE, processed_at = $1 WHERE id = $2`

	if _, err := or.db.Exec(ctx, query, now, id); err != nil {
		return err
	}

	return nil
}
