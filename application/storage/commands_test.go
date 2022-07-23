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
	"github.com/kammeph/school-book-storage-service/infrastructure/memory"
	"github.com/stretchr/testify/assert"
)

func newMemoryStoreWithDefaultEvents() *memory.MemoryStore {
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
	return memory.NewMemoryStoreWithEvents([]domain_common.Event{&eventForRemove, &eventForUpdate})
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
