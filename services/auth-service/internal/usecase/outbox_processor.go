package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type OutboxReader interface {
	GetUnprocessedEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error)
}

type OutboxWriter interface {
	MarkAsProcessed(ctx context.Context, id string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload map[string]interface{}) error
}

type OutboxProcessor struct {
	reader    OutboxReader
	writer    OutboxWriter
	publisher EventPublisher
	logger    *slog.Logger
	interval  time.Duration
	batchSize int
}

func NewOutboxProcessor(
	reader OutboxReader,
	writer OutboxWriter,
	publisher EventPublisher,
	logger *slog.Logger,
	interval time.Duration,
	batchSize int,
) *OutboxProcessor {
	return &OutboxProcessor{
		reader:    reader,
		writer:    writer,
		publisher: publisher,
		logger:    logger,
		interval:  interval,
		batchSize: batchSize,
	}
}

func (op *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(op.interval)
	defer ticker.Stop()

	op.logger.Info("outbox processor started")

	for {
		select {
		case <-ctx.Done():
			op.logger.Info("outbox processor stopped")
			return
		case <-ticker.C:
			op.processBatch(ctx)
		}
	}
}

func (op *OutboxProcessor) processBatch(ctx context.Context) {
	events, err := op.reader.GetUnprocessedEvents(ctx, op.batchSize)
	if err != nil {
		op.logger.Error("failed to get unprocessed events", slog.Any("err", err))
		return
	}

	if len(events) == 0 {
		return
	}

	op.logger.Info("processing outbox events", slog.Int("count", len(events)))

	for _, event := range events {
		if err := op.publisher.Publish(ctx, event.EventType, event.Payload); err != nil {
			op.logger.Error("failed to publish event",
				slog.String("event_id", event.ID.String()),
				slog.String("event_type", event.EventType),
				slog.Any("err", err))
			continue
		}

		if err := op.writer.MarkAsProcessed(ctx, event.ID.String()); err != nil {
			op.logger.Error("failed to mark event as processed",
				slog.String("event_id", event.ID.String()),
				slog.Any("err", err))
		}
	}
}
