package common

import "github.com/google/uuid"

type Command interface {
	AggregateID() uuid.UUID
}

type CommandModel struct {
	ID uuid.UUID
}

func (c CommandModel) AggregateID() uuid.UUID {
	return c.ID
}
