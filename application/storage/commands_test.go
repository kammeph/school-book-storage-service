package storage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/application/storage"
	domain_common "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/stretchr/testify/assert"
)

type EntityAggregate struct {
	domain_common.AggregateModel
}

func (e *EntityAggregate) On(event domain_common.Event) error {
	return nil
}

type memoryStore struct {
	eventsById map[string][]domain_common.Event
}

func (s *memoryStore) Save(ctx context.Context, aggregate domain_common.Aggregate) error {
	if _, ok := s.eventsById[aggregate.AggregateID()]; !ok {
		s.eventsById[aggregate.AggregateID()] = []domain_common.Event{}
	}
	history := append(s.eventsById[aggregate.AggregateID()], aggregate.DomainEvents()...)
	s.eventsById[aggregate.AggregateID()] = history
	return nil
}

func (s *memoryStore) Load(ctx context.Context, aggregate domain_common.Aggregate) error {
	events, ok := s.eventsById[aggregate.AggregateID()]
	if !ok {
		return nil
	}
	if err := aggregate.Load(events); err != nil {
		return err
	}
	return nil
}

var store = &memoryStore{eventsById: map[string][]domain_common.Event{}}

func TestHandleAddStorage(t *testing.T) {
	handler := storage.NewAddStorageCommandHandler(store, nil)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorageCommand{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	storageID, err := handler.Handle(ctx, add)
	assert.Nil(t, err)
	assert.NotZero(t, storageID, 3)
}

func TestHandleRemoveStorage(t *testing.T) {
	addHandler := storage.NewAddStorageCommandHandler(store, nil)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorageCommand{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	storageID, err := addHandler.Handle(ctx, add)
	removeHandler := storage.NewRemoveStorageCommandHandler(store, nil)
	remove := storage.RemoveStorageCommand{CommandModel: common.CommandModel{ID: commandId}, StorageID: storageID, Reason: "test"}
	err = removeHandler.Handle(ctx, remove)
	assert.Nil(t, err)
}

func TestHandleSetStorageName(t *testing.T) {
	addHandler := storage.NewAddStorageCommandHandler(store, nil)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorageCommand{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	storageID, err := addHandler.Handle(ctx, add)
	setNameHandler := storage.NewRenameStorageCommandHandler(store, nil)
	setName := storage.RenameStorageCommand{
		CommandModel: common.CommandModel{ID: commandId},
		StorageID:    storageID,
		Name:         "storage name set",
		Reason:       "test",
	}
	err = setNameHandler.Handle(ctx, setName)
	assert.Nil(t, err)
}

func TestHandleSetStorageLocation(t *testing.T) {
	addHandler := storage.NewAddStorageCommandHandler(store, nil)
	commandId := uuid.New().String()
	ctx := context.Background()
	add := storage.AddStorageCommand{CommandModel: common.CommandModel{ID: commandId}, Name: "storage", Location: "location"}
	dto, err := addHandler.Handle(ctx, add)
	setLocationHandler := storage.NewRelocateStorageCommandHandler(store, nil)
	setLocation := storage.RelocateStorageCommand{
		CommandModel: common.CommandModel{ID: commandId},
		StorageID:    dto,
		Location:     "location set",
		Reason:       "test",
	}
	err = setLocationHandler.Handle(ctx, setLocation)
	assert.Nil(t, err)
}
