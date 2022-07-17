package common

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

// History represents
type History []common.Event

// Len implements sort.Interface
func (h History) Len() int {
	return len(h)
}

// Swap implements sort.Interface
func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Less implements sort.Interface
func (h History) Less(i, j int) bool {
	return h[i].EventVersion() < h[j].EventVersion()
}

type Store interface {
	Load(ctx context.Context, aggregate common.Aggregate) error
	Save(ctx context.Context, aggregate common.Aggregate) error
}
