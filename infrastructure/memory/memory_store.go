package memory

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

type MemoryStore struct {
	eventsById map[string][]common.Event
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{eventsById: map[string][]common.Event{}}
}

func (s *MemoryStore) Save(ctx context.Context, aggregateID string, records ...common.Event) error {
	if _, ok := s.eventsById[aggregateID]; !ok {
		s.eventsById[aggregateID] = []common.Event{}
	}
	history := append(s.eventsById[aggregateID], records...)
	s.eventsById[aggregateID] = history
	return nil
}

func (s *MemoryStore) Load(ctx context.Context, aggregateID string) ([]common.Event, error) {
	_, ok := s.eventsById[aggregateID]
	if ok {
		return s.eventsById[aggregateID], nil
	}
	return nil, nil
}
