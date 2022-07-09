package storage

import (
	"time"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

type StorageCreated struct {
	common.EventModel
	StorageID string
}

func NewStorageCreated(aggregate SchoolAggregateRoot, storageID string) *StorageCreated {
	return &StorageCreated{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now()},
		StorageID: storageID,
	}
}

type StorageRemoved struct {
	common.EventModel
	StorageID string
	Reason    string
}

func NewStorageRemoved(aggregate SchoolAggregateRoot, storageID string, reason string) *StorageRemoved {
	return &StorageRemoved{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now()},
		StorageID: storageID,
		Reason:    reason,
	}
}

type StorageNameSet struct {
	common.EventModel
	StorageID string
	Name      string
	Reason    string
}

func NewStorageNameSet(aggregate SchoolAggregateRoot, storageID string, name, reason string) *StorageNameSet {
	return &StorageNameSet{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now()},
		StorageID: storageID,
		Name:      name,
		Reason:    reason,
	}
}

type StorageLocationSet struct {
	common.EventModel
	StorageID string
	Location  string
	Reason    string
}

func NewStorageLocationSet(aggregate SchoolAggregateRoot, storageID string, location, reason string) *StorageLocationSet {
	return &StorageLocationSet{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now()},
		StorageID: storageID,
		Location:  location,
		Reason:    reason,
	}
}
