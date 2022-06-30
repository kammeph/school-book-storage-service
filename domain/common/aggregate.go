package common

import "github.com/google/uuid"

type AggregateModel struct {
	ID      uuid.UUID
	Version int
}

func (a AggregateModel) AggregateID() uuid.UUID {
	return a.ID
}

func (a AggregateModel) AggregateVersion() int {
	return a.Version
}

type Aggregate interface {
	AggregateID() uuid.UUID
	AggregateVersion() int
	On(event Event) error
}
