package common

import "github.com/google/uuid"

type AggregateModel struct {
	ID      uuid.UUID
	Version int
	Events  []Event
}

func (a AggregateModel) AggregateID() uuid.UUID {
	return a.ID
}

func (a AggregateModel) AggregateVersion() int {
	return a.Version
}

func (a AggregateModel) DomainEvents() []Event {
	return a.Events
}

type Aggregate interface {
	AggregateID() uuid.UUID
	AggregateVersion() int
	DomainEvents() []Event
	On(event Event) error
}
