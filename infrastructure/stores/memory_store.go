package stores

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application/common"
)

type MemoryStore struct {
	eventsById map[string][]common.Record
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{eventsById: map[string][]common.Record{}}
}

func (s *MemoryStore) Save(ctx context.Context, aggregateID string, records ...common.Record) error {
	if _, ok := s.eventsById[aggregateID]; !ok {
		s.eventsById[aggregateID] = []common.Record{}
	}
	history := append(s.eventsById[aggregateID], records...)
	s.eventsById[aggregateID] = history
	return nil
}

func (s *MemoryStore) Load(ctx context.Context, aggregateID string) ([]common.Record, error) {
	_, ok := s.eventsById[aggregateID]
	if ok {
		return s.eventsById[aggregateID], nil
	}
	return nil, nil
}
