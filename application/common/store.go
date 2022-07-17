package common

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Store interface {
	Load(ctx context.Context, aggregate common.Aggregate) error
	Save(ctx context.Context, aggregate common.Aggregate) error
}
