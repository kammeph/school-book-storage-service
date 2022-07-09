package common

import (
	"time"
)

type Event interface {
	AggregateID() string
	EventVersion() int
	EventAt() time.Time
}

type EventTyper interface {
	EventType() string
}

type EventModel struct {
	ID      string
	Version int
	At      time.Time
}

func (m EventModel) AggregateID() string {
	return m.ID
}

func (m EventModel) EventVersion() int {
	return m.Version
}

func (m EventModel) EventAt() time.Time {
	return m.At
}
