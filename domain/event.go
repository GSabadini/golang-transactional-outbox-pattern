package domain

import (
	"context"
	"time"
)

type (
	Producer interface {
		Publish(context.Context, Event) error
		PublishBatch(context.Context, []Event) error
	}

	Event struct {
		Domain    string    `json:"domain"`
		Type      string    `json:"type"`
		Body      string    `json:"body"`
		Timestamp time.Time `json:"timestamp"`
	}
)

func NewEvent(domain string, eventType string, body string, timestamp time.Time) Event {
	return Event{
		Domain:    domain,
		Type:      eventType,
		Body:      body,
		Timestamp: timestamp,
	}
}
