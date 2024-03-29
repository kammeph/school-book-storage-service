package storageapp_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/application/storageapp"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/stretchr/testify/assert"
)

func newMemoryStoreWithDefaultEvents() *memory.MemoryStore {
	eventDataForRemove, _ := json.Marshal(storagedomain.StorageAddedEvent{
		SchoolID:  "school",
		StorageID: "testRemove",
		Name:      "Closet to Remove",
		Location:  "Room 12",
	})
	eventForRemove := domain.EventModel{
		ID:      "school",
		Type:    storagedomain.StorageAdded,
		Version: 1,
		At:      time.Now(),
		Data:    string(eventDataForRemove),
	}
	eventDataForUpdate, _ := json.Marshal(storagedomain.StorageAddedEvent{
		SchoolID:  "school",
		StorageID: "testUpdate",
		Name:      "Closet to Update",
		Location:  "Room 12",
	})
	eventForUpdate := domain.EventModel{
		ID:      "school",
		Type:    storagedomain.StorageAdded,
		Version: 2,
		At:      time.Now(),
		Data:    string(eventDataForUpdate),
	}
	return memory.NewMemoryStoreWithEvents([]domain.Event{&eventForRemove, &eventForUpdate})
}

var store = newMemoryStoreWithDefaultEvents()

func TestHandleAddStorage(t *testing.T) {
	ctx := context.Background()
	handler := storageapp.NewAddStorageCommandHandler(store, nil)
	command := storageapp.AddStorageCommand{CommandModel: application.CommandModel{ID: "school"}, Name: "storage", Location: "location"}
	storageID, err := handler.Handle(ctx, command)
	assert.Nil(t, err)
	assert.NotZero(t, storageID, 3)
}

func TestHandleRemoveStorage(t *testing.T) {
	ctx := context.Background()
	removeHandler := storageapp.NewRemoveStorageCommandHandler(store, nil)
	remove := storageapp.RemoveStorageCommand{CommandModel: application.CommandModel{ID: "school"}, StorageID: "testRemove", Reason: "test"}
	err := removeHandler.Handle(ctx, remove)
	assert.Nil(t, err)
}

func TestHandleSetStorageName(t *testing.T) {
	ctx := context.Background()
	handler := storageapp.NewRenameStorageCommandHandler(store, nil)
	command := storageapp.RenameStorageCommand{
		CommandModel: application.CommandModel{ID: "school"},
		StorageID:    "testUpdate",
		Name:         "storage name set",
		Reason:       "test",
	}
	err := handler.Handle(ctx, command)
	assert.Nil(t, err)
}

func TestHandleSetStorageLocation(t *testing.T) {
	ctx := context.Background()
	handler := storageapp.NewRelocateStorageCommandHandler(store, nil)
	command := storageapp.RelocateStorageCommand{
		CommandModel: application.CommandModel{ID: "school"},
		StorageID:    "testUpdate",
		Location:     "location set",
		Reason:       "test",
	}
	err := handler.Handle(ctx, command)
	assert.Nil(t, err)
}
