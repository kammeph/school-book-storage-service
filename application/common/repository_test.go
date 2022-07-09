package common_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/serializers"
	"github.com/stretchr/testify/assert"
)

type EntityAggregate struct {
	domain.AggregateModel
	Entity Entity
}

type Entity struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EntityCreated struct {
	domain.EventModel
	EntityID string
}

type EntityNameSet struct {
	domain.EventModel
	Name string
}

func (e *EntityAggregate) On(event domain.Event) error {
	switch v := event.(type) {
	case *EntityCreated:
		e.ID = v.AggregateID()
		e.Version = v.EventVersion()
		e.Entity.ID = v.EntityID
		e.Entity.CreatedAt = v.EventAt()
	case *EntityNameSet:
		e.Version = v.EventVersion()
		e.Entity.Name = v.Name
		e.Entity.UpdatedAt = v.EventAt()
	default:
		return fmt.Errorf("Unknown event type: %t", event)
	}

	return nil
}

type CreateEntity struct {
	common.CommandModel
}

type SetEntityName struct {
	common.CommandModel
	Name string
}

type memoryStore struct {
	eventsById map[string][]common.Record
}

func (s *memoryStore) Save(ctx context.Context, aggregateID string, records ...common.Record) error {
	if _, ok := s.eventsById[aggregateID]; !ok {
		s.eventsById[aggregateID] = []common.Record{}
	}
	history := append(s.eventsById[aggregateID], records...)
	s.eventsById[aggregateID] = history
	return nil
}

func (s *memoryStore) Load(ctx context.Context, aggregateID string) ([]common.Record, error) {
	_, ok := s.eventsById[aggregateID]
	if ok {
		return s.eventsById[aggregateID], nil
	}
	return nil, nil
}

func TestNew(t *testing.T) {
	repository := common.NewRepository(
		&EntityAggregate{},
		&memoryStore{eventsById: map[string][]common.Record{}},
		serializers.NewJSONSerializerWithEvents())
	assert.NotNil(t, repository)
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	aggregateID := uuid.New().String()
	repository := common.NewRepository(
		&EntityAggregate{},
		&memoryStore{eventsById: map[string][]common.Record{}},
		serializers.NewJSONSerializerWithEvents(
			EntityCreated{},
			EntityNameSet{}))
	aggregate, err := repository.Load(ctx, aggregateID)
	assert.Nil(t, err)
	assert.NotNil(t, aggregate)
	assert.Equal(t, aggregate.AggregateID(), aggregateID)
	_, ok := aggregate.(*EntityAggregate)
	assert.True(t, ok)
}

func TestSave(t *testing.T) {
	ctx := context.Background()
	aggregateID := uuid.New().String()
	repository := common.NewRepository(
		&EntityAggregate{},
		&memoryStore{eventsById: map[string][]common.Record{}},
		serializers.NewJSONSerializerWithEvents(
			EntityCreated{},
			EntityNameSet{}))
	aggregate, err := repository.Load(ctx, aggregateID)
	assert.Nil(t, err)
	assert.NotNil(t, aggregate)
	a, ok := aggregate.(*EntityAggregate)
	assert.True(t, ok)
	entityID := uuid.New().String()
	createdEvent := EntityCreated{EventModel: domain.EventModel{ID: aggregateID, Version: 1, At: time.Now()}, EntityID: entityID}
	a.Events = append(a.Events, createdEvent)
	nameSetEvent := EntityNameSet{EventModel: domain.EventModel{ID: aggregateID, Version: 2, At: time.Now()}, Name: "entity"}
	a.Events = append(a.Events, nameSetEvent)
	repository.Save(ctx, a)
	aggregate, err = repository.Load(ctx, aggregateID)
	assert.Nil(t, err)
	assert.NotNil(t, aggregate)
	a, ok = aggregate.(*EntityAggregate)
	assert.True(t, ok)
	assert.Equal(t, a.AggregateID(), aggregateID)
	assert.Equal(t, a.Entity.ID, entityID)
	assert.Equal(t, a.Entity.Name, "entity")
}
