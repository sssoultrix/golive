package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type EventPublisher struct {
	client      *http.Client
	profileURL  string
	logger      *slog.Logger
}

func NewEventPublisher(profileURL string, logger *slog.Logger) *EventPublisher {
	return &EventPublisher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		profileURL: profileURL,
		logger:     logger,
	}
}

func (ep *EventPublisher) Publish(ctx context.Context, eventType string, payload map[string]interface{}) error {
	switch eventType {
	case "UserRegistered":
		return ep.publishUserRegistered(ctx, payload)
	default:
		ep.logger.Warn("unknown event type", slog.String("event_type", eventType))
		return nil
	}
}

func (ep *EventPublisher) publishUserRegistered(ctx context.Context, payload map[string]interface{}) error {
	userID, ok := payload["user_id"].(string)
	if !ok {
		ep.logger.Error("invalid user_id in payload")
		return fmt.Errorf("invalid user_id in payload")
	}

	login, ok := payload["login"].(string)
	if !ok {
		ep.logger.Error("invalid login in payload")
		return fmt.Errorf("invalid login in payload")
	}

	createProfileReq := map[string]interface{}{
		"name":  login,
		"login": login,
		"email": login + "@example.com",
	}

	body, err := json.Marshal(createProfileReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ep.profileURL+"/profiles", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ep.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		ep.logger.Warn("failed to create profile",
			slog.String("user_id", userID),
			slog.Int("status", resp.StatusCode))
		return nil
	}

	ep.logger.Info("profile created successfully", slog.String("user_id", userID))
	return nil
}
