package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain"
)

type StorageEventHandler struct {
	repository StorageWithBooksRepository
}

func NewStorageEventHandler(repository StorageWithBooksRepository) EventHandler {
	return &StorageEventHandler{repository}
}

func (h StorageEventHandler) Handle(ctx context.Context, eventBytes []byte) {
	event := &domain.EventModel{}
	if err := json.Unmarshal(eventBytes, event); err != nil {
		fmt.Print(err)
	}
	switch event.EventType() {
	case domain.StorageAdded:
		h.handleStorageAdded(ctx, event)
		return
	case domain.StorageRemoved:
		h.handleStorageRemoved(ctx, event)
		return
	case domain.StorageRenamed:
		h.handleStorageRenamed(ctx, event)
		return
	case domain.StorageRelocated:
		h.handleStorageRelocated(ctx, event)
		return
	default:
		return
	}
}

func (h StorageEventHandler) handleStorageAdded(ctx context.Context, event domain.Event) {
	storageAdded := domain.StorageAddedEvent{}
	if err := json.Unmarshal([]byte(event.EventData()), &storageAdded); err != nil {
		fmt.Print(err)
	}
	storage := domain.NewStorageWithBooks(
		storageAdded.SchoolID,
		storageAdded.StorageID,
		storageAdded.Name,
		storageAdded.Location)
	h.repository.InsertStorage(ctx, storage)
}

func (h StorageEventHandler) handleStorageRemoved(ctx context.Context, event domain.Event) {
	storageRemoved := domain.StorageRemovedEvent{}
	err := json.Unmarshal([]byte(event.EventData()), &storageRemoved)
	if err != nil {
		fmt.Print(err)
	}
	h.repository.DeleteStorage(ctx, storageRemoved.StorageID)
}

func (h StorageEventHandler) handleStorageRenamed(ctx context.Context, event domain.Event) {
	storageRenamed := domain.StorageRenamedEvent{}
	err := json.Unmarshal([]byte(event.EventData()), &storageRenamed)
	if err != nil {
		fmt.Print(err)
	}
	h.repository.UpdateStorageName(ctx, storageRenamed.StorageID, storageRenamed.Name)
}

func (h StorageEventHandler) handleStorageRelocated(ctx context.Context, event domain.Event) {
	storageRelocated := domain.StorageRelocatedEvent{}
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