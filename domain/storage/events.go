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
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Location  string `json:"location"`
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
	StorageID string `json:"storageId"`
	Reason    string `json:"reason"`
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
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
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
	StorageID string `json:"storageId"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
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
