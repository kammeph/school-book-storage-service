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

func (s *MemoryStore) Save(ctx context.Context, aggregate common.Aggregate) error {
	if _, ok := s.eventsById[aggregate.AggregateID()]; !ok {
		s.eventsById[aggregate.AggregateID()] = []common.Event{}
	}
	history := append(s.eventsById[aggregate.AggregateID()], aggregate.DomainEvents()...)
	s.eventsById[aggregate.AggregateID()] = history
	return nil
}

func (s *MemoryStore) Load(ctx context.Context, aggregate common.Aggregate) error {
	events := s.eventsById[aggregate.AggregateID()]
	for _, event := range events {
		if err := aggregate.On(event); err != nil {
			return err
		}
	}
	return nil
}
