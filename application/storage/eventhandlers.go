package storage

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageCreatedEventHandler struct{}

func NewStorageCreatedEventHandler() StorageCreatedEventHandler {
	return StorageCreatedEventHandler{}
}

func (h StorageCreatedEventHandler) Handle(ctx context.Context, eventData []byte) {
	event := domain.StorageCreated{}
	err := json.Unmarshal(eventData, &event)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println()
	fmt.Printf("Received storage created event %v: ", event)
	fmt.Println()
}

type StorageNameSetEventHandler struct{}

func NewStorageNameSetEventHandler() StorageNameSetEventHandler {
	return StorageNameSetEventHandler{}
}

func (h StorageNameSetEventHandler) Handle(ctx context.Context, eventData []byte) {
	event := domain.StorageNameSet{}
	err := json.Unmarshal(eventData, &event)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println()
	fmt.Printf("Received storage name set event %v: ", event)
	fmt.Println()
}
