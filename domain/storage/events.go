package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type StorageCreated struct {
	common.EventModel
	StorageID uuid.UUID
	SchoolID  uuid.UUID
}

func NewStorageCreated(aggregate SchoolAggregateRoot, storageID uuid.UUID) StorageCreated {
	return StorageCreated{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now()},
		StorageID: storageID,
	}
}

type StorageRemoved struct {
	common.EventModel
	StorageID uuid.UUID
	Reason    string
}

func NewStorageRemoved(aggregate SchoolAggregateRoot, storageID uuid.UUID, reason string) StorageRemoved {
	return StorageRemoved{
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
	StorageID uuid.UUID
	Name      string
	Reason    string
}

func NewStorageNameSet(aggregate SchoolAggregateRoot, storageID uuid.UUID, name, reason string) StorageNameSet {
	return StorageNameSet{
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
	StorageID uuid.UUID
	Location  string
	Reason    string
}

func NewStorageLocationSet(aggregate SchoolAggregateRoot, storageID uuid.UUID, location, reason string) StorageLocationSet {
	return StorageLocationSet{
		EventModel: common.EventModel{
			ID:      aggregate.AggregateID(),
			Version: aggregate.AggregateVersion() + 1,
			At:      time.Now()},
		StorageID: storageID,
		Location:  location,
		Reason:    reason,
	}
}
