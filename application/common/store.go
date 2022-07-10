package common

import (
	"context"
)

type Record struct {
	Version int
	Data    string
}

// History represents
type History []Record

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
	return h[i].Version < h[j].Version
}

type Store interface {
	Load(ctx context.Context, aggregateID string) ([]Record, error)
	Save(ctx context.Context, aggregateID string, records ...Record) error
}
