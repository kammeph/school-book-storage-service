package common

import "github.com/google/uuid"

type QueryModel struct {
	ID uuid.UUID
}

type Query interface {
	AggregateID() uuid.UUID
}

func (q QueryModel) AggregateID() uuid.UUID {
	return q.ID
}
