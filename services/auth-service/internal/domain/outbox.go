package domain

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID         uuid.UUID
	EventType  string
	Payload    map[string]interface{}
	Processed  bool
	CreatedAt  time.Time
	ProcessedAt *time.Time
}

func NewOutboxEvent(eventType string, payload map[string]interface{}) *OutboxEvent {
	return &OutboxEvent{
		ID:        uuid.New(),
		EventType: eventType,
		Payload:   payload,
		Processed: false,
		CreatedAt: time.Now().UTC(),
	}
}
