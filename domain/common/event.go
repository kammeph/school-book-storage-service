package common

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	AggregateID() uuid.UUID
	EventVersion() int
	EventAt() time.Time
}

type EventTyper interface {
	EventType() string
}

type EventModel struct {
	ID      uuid.UUID
	Version int
	At      time.Time
}

func (m EventModel) AggregateID() uuid.UUID {
	return m.ID
}

func (m EventModel) EventVersion() int {
	return m.Version
}

func (m EventModel) EventAt() time.Time {
	return m.At
}
