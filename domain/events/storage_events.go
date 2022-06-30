package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type StorageCreated struct {
	common.EventModel
	StorageID uuid.UUID
}

func NewStorageCreated(aggregateID uuid.UUID, version int, storageID uuid.UUID) StorageCreated {
	return StorageCreated{
		EventModel: common.EventModel{
			ID:      aggregateID,
			Version: version,
			At:      time.Now()},
		StorageID: storageID}
}

type StorageRemoved struct {
	common.EventModel
	StorageID uuid.UUID
	Reason    string
}

func NewStorageRemoved(aggregateID uuid.UUID, version int, storageID uuid.UUID, reason string) StorageRemoved {
	return StorageRemoved{
		EventModel: common.EventModel{
			ID:      aggregateID,
			Version: version,
			At:      time.Now()},
		StorageID: storageID,
		Reason:    reason}
}

type StorageNameSet struct {
	common.EventModel
	StorageID uuid.UUID
	Name      string
	Reason    string
}

func NewStorageNameSet(aggregateID uuid.UUID, version int, storageID uuid.UUID, name string, reason string) StorageNameSet {
	return StorageNameSet{
		EventModel: common.EventModel{
			ID:      aggregateID,
			Version: version,
			At:      time.Now()},
		StorageID: storageID,
		Name:      name,
		Reason:    reason}
}

type StorageLocationSet struct {
	common.EventModel
	StorageID uuid.UUID
	Location  string
	Reason    string
}

func NewStorageLocationSet(aggregateID uuid.UUID, version int, storageID uuid.UUID, location string, reason string) StorageLocationSet {
	return StorageLocationSet{
		EventModel: common.EventModel{
			ID:      aggregateID,
			Version: version,
			At:      time.Now()},
		StorageID: storageID,
		Location:  location,
		Reason:    reason}
}
