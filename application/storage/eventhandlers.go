package storage

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageAddedEventHandler struct{}

func NewStorageAddedEventHandler() StorageAddedEventHandler {
	return StorageAddedEventHandler{}
}

func (h StorageAddedEventHandler) Handle(ctx context.Context, eventData []byte) {
	event := domain.StorageAdded{}
	err := json.Unmarshal(eventData, &event)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println()
	fmt.Printf("Received storage added event %s: ", eventData)
	fmt.Println()
}

type StorageRenamedEventHandler struct{}

func NewStorageRenamedSetEventHandler() StorageRenamedEventHandler {
	return StorageRenamedEventHandler{}
}

func (h StorageRenamedEventHandler) Handle(ctx context.Context, eventData []byte) {
	event := domain.StorageRenamed{}
	err := json.Unmarshal(eventData, &event)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println()
	fmt.Printf("Received storage renamed event %s: ", eventData)
	fmt.Println()
}
