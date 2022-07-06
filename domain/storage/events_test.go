package storage_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewStorageCreated(t *testing.T) {
	storgeID, aggregate := newTestAggregate()
	storageCreated := storage.NewStorageCreated(aggregate, storgeID)
	assert.NotNil(t, storageCreated)
	assert.Equal(t, storageCreated.AggregateID(), aggregate.AggregateID())
	assert.Equal(t, storageCreated.Version, aggregate.AggregateVersion()+1)
	assert.NotZero(t, storageCreated.EventAt())
	assert.Equal(t, storageCreated.StorageID, storgeID)
}

func TestNewStorageRemoved(t *testing.T) {
	storgeID, aggregate := newTestAggregate()
	storageRemoved := storage.NewStorageRemoved(aggregate, storgeID, "test")
	assert.NotNil(t, storageRemoved)
	assert.Equal(t, storageRemoved.AggregateID(), aggregate.AggregateID())
	assert.Equal(t, storageRemoved.Version, aggregate.AggregateVersion()+1)
	assert.NotZero(t, storageRemoved.EventAt())
	assert.Equal(t, storageRemoved.StorageID, storgeID)
	assert.Equal(t, storageRemoved.Reason, "test")
}

func TestNewStorageNameSet(t *testing.T) {
	storgeID, aggregate := newTestAggregate()
	storageNameSet := storage.NewStorageNameSet(aggregate, storgeID, "storage", "test")
	assert.NotNil(t, storageNameSet)
	assert.Equal(t, storageNameSet.AggregateID(), aggregate.AggregateID())
	assert.Equal(t, storageNameSet.Version, aggregate.AggregateVersion()+1)
	assert.NotZero(t, storageNameSet.EventAt())
	assert.Equal(t, storageNameSet.StorageID, storgeID)
	assert.Equal(t, storageNameSet.Name, "storage")
	assert.Equal(t, storageNameSet.Reason, "test")
}

func TestNewStorageLocationSet(t *testing.T) {
	storgeID, aggregate := newTestAggregate()
	storageLocationSet := storage.NewStorageLocationSet(aggregate, storgeID, "location", "test")
	assert.NotNil(t, storageLocationSet)
	assert.Equal(t, storageLocationSet.AggregateID(), aggregate.AggregateID())
	assert.Equal(t, storageLocationSet.Version, aggregate.AggregateVersion()+1)
	assert.NotZero(t, storageLocationSet.EventAt())
	assert.Equal(t, storageLocationSet.StorageID, storgeID)
	assert.Equal(t, storageLocationSet.Location, "location")
	assert.Equal(t, storageLocationSet.Reason, "test")
}
