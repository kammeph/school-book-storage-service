package common

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Store interface {
	Load(ctx context.Context, aggregateID string) ([]common.Event, error)
	Save(ctx context.Context, events []common.Event) error
}
