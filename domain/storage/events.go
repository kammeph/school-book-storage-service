package storage

import (
	"time"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

type StorageAdded struct {
	common.EventModel
	StorageID string
	Name      string
	Location  string
}

func NewStorageAdded(aggregate StorageAggregateRoot, storageID, name, location string) *StorageAdded {
	return &StorageAdded{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now(),
		},
		StorageID: storageID,
		Name:      name,
		Location:  location,
	}
}

type StorageRemoved struct {
	common.EventModel
	StorageID string
	Reason    string
}

func NewStorageRemoved(aggregate StorageAggregateRoot, storageID string, reason string) *StorageRemoved {
	return &StorageRemoved{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now(),
		},
		StorageID: storageID,
		Reason:    reason,
	}
}

type StorageRenamed struct {
	common.EventModel
	StorageID string
	Name      string
	Reason    string
}

func NewStorageRenamed(aggregate StorageAggregateRoot, storageID string, name, reason string) *StorageRenamed {
	return &StorageRenamed{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now(),
		},
		StorageID: storageID,
		Name:      name,
		Reason:    reason,
	}
}

type StorageRelocated struct {
	common.EventModel
	StorageID string
	Location  string
	Reason    string
}

func NewStorageRelocated(aggregate StorageAggregateRoot, storageID string, location, reason string) *StorageRelocated {
	return &StorageRelocated{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now(),
		},
		StorageID: storageID,
		Location:  location,
		Reason:    reason,
	}
}
