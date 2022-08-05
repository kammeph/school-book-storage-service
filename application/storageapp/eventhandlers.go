package storageapp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
)

type StorageEventHandler struct {
	repository StorageWithBooksRepository
}

func NewStorageEventHandler(repository StorageWithBooksRepository) application.EventHandler {
	return &StorageEventHandler{repository}
}

func (h StorageEventHandler) Handle(ctx context.Context, eventBytes []byte) {
	event := &domain.EventModel{}
	if err := json.Unmarshal(eventBytes, event); err != nil {
		fmt.Print(err)
	}
	switch event.EventType() {
	case storagedomain.StorageAdded:
		h.handleStorageAdded(ctx, event)
		return
	case storagedomain.StorageRemoved:
		h.handleStorageRemoved(ctx, event)
		return
	case storagedomain.StorageRenamed:
		h.handleStorageRenamed(ctx, event)
		return
	case storagedomain.StorageRelocated:
		h.handleStorageRelocated(ctx, event)
		return
	default:
		return
	}
}

func (h StorageEventHandler) handleStorageAdded(ctx context.Context, event domain.Event) {
	storageAdded := storagedomain.StorageAddedEvent{}
	if err := json.Unmarshal([]byte(event.EventData()), &storageAdded); err != nil {
		fmt.Print(err)
	}
	storage := storagedomain.NewStorageWithBooks(
		storageAdded.SchoolID,
		storageAdded.StorageID,
		storageAdded.Name,
		storageAdded.Location)
	h.repository.InsertStorage(ctx, storage)
}

func (h StorageEventHandler) handleStorageRemoved(ctx context.Context, event domain.Event) {
	storageRemoved := storagedomain.StorageRemovedEvent{}
	err := json.Unmarshal([]byte(event.EventData()), &storageRemoved)
	if err != nil {
		fmt.Print(err)
	}
	h.repository.DeleteStorage(ctx, storageRemoved.StorageID)
}

func (h StorageEventHandler) handleStorageRenamed(ctx context.Context, event domain.Event) {
	storageRenamed := storagedomain.StorageRenamedEvent{}
	err := json.Unmarshal([]byte(event.EventData()), &storageRenamed)
	if err != nil {
		fmt.Print(err)
	}
	h.repository.UpdateStorageName(ctx, storageRenamed.StorageID, storageRenamed.Name)
}

func (h StorageEventHandler) handleStorageRelocated(ctx context.Context, event domain.Event) {
	storageRelocated := storagedomain.StorageRelocatedEvent{}
	err := json.Unmarshal([]byte(event.EventData()), &storageRelocated)
	if err != nil {
		fmt.Print(err)
	}
	h.repository.UpdateStorageLocation(ctx, storageRelocated.StorageID, storageRelocated.Location)
}

type TestHandler struct{}

func (h TestHandler) Handle(ctx context.Context, eventBytes []byte) {
	fmt.Printf("%s", eventBytes)
}
