package storage_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/application/storage"
	domain_common "github.com/kammeph/school-book-storage-service/domain/common"
	domain_storage "github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/stretchr/testify/assert"
)

type memoryStore struct {
	eventsById map[string][]domain_common.Event
}

func newMemoryStoreWithDefaultEvents() *memoryStore {
	store := &memoryStore{eventsById: map[string][]domain_common.Event{}}
	eventDataForRemove, _ := json.Marshal(domain_storage.StorageAddedEvent{
		SchoolID:  "school",
		StorageID: "testRemove",
		Name:      "Closet to Remove",
		Location:  "Room 12",
	})
	eventForRemove := domain_common.EventModel{
		ID:      "school",
		Type:    domain_storage.StorageAdded,
		Version: 1,
		At:      time.Now(),
		Data:    string(eventDataForRemove),
	}
	eventDataForUpdate, _ := json.Marshal(domain_storage.StorageAddedEvent{
		SchoolID:  "school",
		StorageID: "testUpdate",
		Name:      "Closet to Update",
		Location:  "Room 12",
	})
	eventForUpdate := domain_common.EventModel{
		ID:      "school",
		Type:    domain_storage.StorageAdded,
		Version: 2,
		At:      time.Now(),
		Data:    string(eventDataForUpdate),
	}
	store.eventsById["school"] = []domain_common.Event{&eventForRemove, &eventForUpdate}
	return store
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

var store = newMemoryStoreWithDefaultEvents()

func TestHandleAddStorage(t *testing.T) {
	ctx := context.Background()
	handler := storage.NewAddStorageCommandHandler(store, nil)
	command := storage.AddStorageCommand{CommandModel: common.CommandModel{ID: "school"}, Name: "storage", Location: "location"}
	storageID, err := handler.Handle(ctx, command)
	assert.Nil(t, err)
	assert.NotZero(t, storageID, 3)
}

func TestHandleRemoveStorage(t *testing.T) {
	ctx := context.Background()
	removeHandler := storage.NewRemoveStorageCommandHandler(store, nil)
	remove := storage.RemoveStorageCommand{CommandModel: common.CommandModel{ID: "school"}, StorageID: "testRemove", Reason: "test"}
	err := removeHandler.Handle(ctx, remove)
	assert.Nil(t, err)
}

func TestHandleSetStorageName(t *testing.T) {
	ctx := context.Background()
	handler := storage.NewRenameStorageCommandHandler(store, nil)
	command := storage.RenameStorageCommand{
		CommandModel: common.CommandModel{ID: "school"},
		StorageID:    "testUpdate",
		Name:         "storage name set",
		Reason:       "test",
	}
	err := handler.Handle(ctx, command)
	assert.Nil(t, err)
}

func TestHandleSetStorageLocation(t *testing.T) {
	ctx := context.Background()
	handler := storage.NewRelocateStorageCommandHandler(store, nil)
	command := storage.RelocateStorageCommand{
		CommandModel: common.CommandModel{ID: "school"},
		StorageID:    "testUpdate",
		Location:     "location set",
		Reason:       "test",
	}
	err := handler.Handle(ctx, command)
	assert.Nil(t, err)
}
