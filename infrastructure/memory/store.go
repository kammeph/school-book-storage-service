package memory

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
)

type MemoryStore struct {
	eventsById map[string][]domain.Event
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{eventsById: map[string][]domain.Event{}}
}

func NewMemoryStoreWithEvents(events []domain.Event) *MemoryStore {
	store := NewMemoryStore()
	if len(events) == 0 {
		return store
	}
	store.Save(context.TODO(), events)
	return store
}

func (s *MemoryStore) Save(ctx context.Context, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	aggregateID := events[0].AggregateID()
	if _, ok := s.eventsById[aggregateID]; !ok {
		s.eventsById[aggregateID] = []domain.Event{}
	}
	history := append(s.eventsById[aggregateID], events...)
	s.eventsById[aggregateID] = history
	return nil
}

func (s *MemoryStore) Load(ctx context.Context, aggregateID string) ([]domain.Event, error) {
	events, ok := s.eventsById[aggregateID]
	if !ok {
		return nil, nil
	}
	return events, nil
}
