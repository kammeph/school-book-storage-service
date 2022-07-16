package storage

import (
	"github.com/kammeph/school-book-storage-service/domain/common"
)

var (
	StorageAdded     = "STORAGE_ADDED"
	StorageRemoved   = "STORAGE_REMOVED"
	StorageRenamed   = "STORAGE_RENAMED"
	StorageRelocated = "STORAGE_RELOCATED"
)

type StorageAddedEvent struct {
	StorageID string
	Name      string
	Location  string
}

func NewStorageAdded(aggregate *StorageAggregate, storageID, name, location string) (common.Event, error) {
	eventData := StorageAddedEvent{
		StorageID: storageID,
		Name:      name,
		Location:  location,
	}
	event := common.NewEvent(aggregate, StorageAdded)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type StorageRemovedEvent struct {
	StorageID string
	Reason    string
}

func NewStorageRemoved(aggregate *StorageAggregate, storageID string, reason string) (common.Event, error) {
	eventData := StorageRemovedEvent{
		StorageID: storageID,
		Reason:    reason,
	}
	event := common.NewEvent(aggregate, StorageRemoved)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type StorageRenamedEvent struct {
	StorageID string
	Name      string
	Reason    string
}

func NewStorageRenamed(aggregate *StorageAggregate, storageID string, name, reason string) (common.Event, error) {
	eventData := StorageRenamedEvent{
		StorageID: storageID,
		Name:      name,
		Reason:    reason,
	}
	event := common.NewEvent(aggregate, StorageRenamed)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type StorageRelocatedEvent struct {
	StorageID string
	Location  string
	Reason    string
}

func NewStorageRelocated(aggregate *StorageAggregate, storageID string, location, reason string) (common.Event, error) {
	eventData := StorageRelocatedEvent{
		StorageID: storageID,
		Location:  location,
		Reason:    reason,
	}
	event := common.NewEvent(aggregate, StorageRelocated)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}
