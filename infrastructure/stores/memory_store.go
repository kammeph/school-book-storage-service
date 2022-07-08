package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
)

type MemoryStore struct {
	eventsById map[string][]common.Record
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{eventsById: map[string][]common.Record{}}
}

func (s *MemoryStore) Save(ctx context.Context, aggregateID uuid.UUID, records ...common.Record) error {
	if _, ok := s.eventsById[aggregateID.String()]; !ok {
		s.eventsById[aggregateID.String()] = []common.Record{}
	}
	history := append(s.eventsById[aggregateID.String()], records...)
	s.eventsById[aggregateID.String()] = history
	return nil
}

func (s *MemoryStore) Load(ctx context.Context, aggregateID uuid.UUID) ([]common.Record, error) {
	_, ok := s.eventsById[aggregateID.String()]
	if ok {
		return s.eventsById[aggregateID.String()], nil
	}
	return nil, nil
}
