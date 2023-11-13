package domain

import "context"

type Producer interface {
	Publish(context.Context, Event) error
}

type Event struct {
	Type      string
	Body      string
	Timestamp string
}

func NewEvent(Type string, body string, timestamp string) Event {
	return Event{Type: Type, Body: body, Timestamp: timestamp}
}
