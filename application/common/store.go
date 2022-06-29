package common

import (
	"context"

	"github.com/google/uuid"
)

type Record struct {
	Version int
	Data    []byte
}

type Store interface {
	Load(ctx context.Context, aggregateID uuid.UUID) ([]Record, error)
	Save(ctx context.Context, aggregateID uuid.UUID, records ...Record) error
}
