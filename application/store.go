package application

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
)

type Store interface {
	Load(ctx context.Context, aggregateID string) ([]domain.Event, error)
	Save(ctx context.Context, events []domain.Event) error
}
