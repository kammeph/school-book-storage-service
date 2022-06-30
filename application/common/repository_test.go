package common_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/serializers"
	"github.com/stretchr/testify/assert"
)

type Entity struct {
	ID        uuid.UUID
	Version   int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EntityCreated struct {
	domain.EventModel
}

type EntityNameSet struct {
	domain.EventModel
	Name string
}

func (e *Entity) CreateEntity(id uuid.UUID) ([]domain.Event, error) {
	return []domain.Event{
			EntityCreated{
				EventModel: domain.EventModel{
					ID:      id,
					Version: e.Version + 1,
					At:      time.Now()}}},
		nil
}

func (e *Entity) SetEntityName(name string) ([]domain.Event, error) {
	return []domain.Event{
			EntityNameSet{
				EventModel: domain.EventModel{
					ID:      e.ID,
					Version: e.Version + 1,
					At:      time.Now()}, Name: name}},
		nil
}

func (e *Entity) On(event domain.Event) error {
	switch v := event.(type) {
	case *EntityCreated:
		e.ID = v.AggregateID()
		e.Version = v.EventVersion()
		e.CreatedAt = v.EventAt()
		e.UpdatedAt = v.EventAt()
	case *EntityNameSet:
		e.Version = v.EventVersion()
		e.UpdatedAt = v.EventAt()
		e.Name = v.Name
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

type EntityCommandHandler struct {
}

func (h *EntityCommandHandler) Apply(ctx context.Context, aggregate domain.Aggregate, command common.Command) ([]domain.Event, error) {
	entity, ok := aggregate.(*Entity)
	if !ok {
		return nil, fmt.Errorf("Incorrect type for aggregate: %t", aggregate)
	}
	switch c := command.(type) {
	case CreateEntity:
		return entity.CreateEntity(c.AggregateID())
	case SetEntityName:
		return entity.SetEntityName(c.Name)
	default:
		return nil, errors.New("Unapplied command")
	}
}

type memoryStore struct {
	eventsById map[string][]common.Record
}

func (s *memoryStore) Save(ctx context.Context, aggregateID uuid.UUID, records ...common.Record) error {
	if _, ok := s.eventsById[aggregateID.String()]; !ok {
		s.eventsById[aggregateID.String()] = []common.Record{}
	}
	history := append(s.eventsById[aggregateID.String()], records...)
	s.eventsById[aggregateID.String()] = history
	return nil
}

func (s *memoryStore) Load(ctx context.Context, aggregateID uuid.UUID) ([]common.Record, error) {
	_, ok := s.eventsById[aggregateID.String()]
	if ok {
		return s.eventsById[aggregateID.String()], nil
	}
	return nil, nil
}

func TestNew(t *testing.T) {
	repository := common.NewRepository(&Entity{}, &memoryStore{eventsById: map[string][]common.Record{}}, serializers.NewJSONSerializer(), &EntityCommandHandler{})
	aggregate := repository.NewAggregate()
	assert.NotNil(t, repository)
	assert.Equal(t, aggregate, &Entity{})
}

func TestSave(t *testing.T) {
	ctx := context.Background()
	entityID := uuid.New()
	repository := common.NewRepository(&Entity{}, &memoryStore{eventsById: map[string][]common.Record{}}, serializers.NewJSONSerializer(EntityCreated{}, EntityNameSet{}), &EntityCommandHandler{})
	createEntity := CreateEntity{CommandModel: common.CommandModel{ID: entityID}}
	setEntityName := SetEntityName{CommandModel: common.CommandModel{ID: entityID}, Name: "Test"}
	err := repository.Save(ctx, createEntity)
	assert.Nil(t, err)
	aggregate, err := repository.Load(ctx, entityID)
	entityAfterCreate, ok := aggregate.(*Entity)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.NotNil(t, aggregate)
	assert.Equal(t, entityAfterCreate.ID, entityID)
	assert.Equal(t, entityAfterCreate.Version, 1)
	assert.NotZero(t, entityAfterCreate.CreatedAt)
	err = repository.Save(ctx, setEntityName)
	newaggregate, err := repository.Load(ctx, entityID)
	entityAfterNameSet, ok := newaggregate.(*Entity)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, entityAfterNameSet.Version, 2)
	assert.NotZero(t, entityAfterNameSet.UpdatedAt)
	assert.Equal(t, entityAfterNameSet.Name, "Test")
}
