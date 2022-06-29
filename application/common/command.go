package common

import (
	"context"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Command interface {
	AggregateID() uuid.UUID
}

type CommandModel struct {
	ID uuid.UUID
}

func (m CommandModel) AggregateID() uuid.UUID {
	return m.ID
}

type CommandHandler interface {
	Apply(ctx context.Context, aggregate common.Aggregate, command Command) ([]common.Event, error)
}
