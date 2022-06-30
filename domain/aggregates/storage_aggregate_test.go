package aggregates_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/aggregates"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/entities"
	"github.com/kammeph/school-book-storage-service/domain/events"
	"github.com/stretchr/testify/assert"
)

type UnknownEvent struct {
	common.EventModel
}

func newTestAggregate(removed bool) (uuid.UUID, aggregates.StorageAggregate) {
	storageID := uuid.New()
	aggregate := aggregates.NewStorageAggregate()
	aggregate.Storages = append(
		aggregate.Storages,
		entities.Storage{
			ID:       storageID,
			Name:     "storage",
			Location: "location",
			Removed:  removed})
	return storageID, aggregate
}

func TestAddStorage(t *testing.T) {
	aggregate := aggregates.NewStorageAggregate()
	createdEvents, err := aggregate.AddStorage("storage", "location")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 3)
	created, createdOk := createdEvents[0].(events.StorageCreated)
	assert.True(t, createdOk)
	assert.Equal(t, created.Version, 1)
	nameSet, nameSetOk := createdEvents[1].(events.StorageNameSet)
	assert.True(t, nameSetOk)
	assert.Equal(t, nameSet.Version, 2)
	assert.Equal(t, nameSet.Name, "storage")
	locationSet, locationOk := createdEvents[2].(events.StorageLocationSet)
	assert.True(t, locationOk)
	assert.Equal(t, locationSet.Version, 3)
	assert.Equal(t, locationSet.Location, "location")
}

func TestAddStorageWithoutName(t *testing.T) {
	aggregate := aggregates.NewStorageAggregate()
	_, err := aggregate.AddStorage("", "location")
	assert.NotNil(t, err)
}

func TestAddStorageWithoutLocation(t *testing.T) {
	aggregate := aggregates.NewStorageAggregate()
	_, err := aggregate.AddStorage("storage", "")
	assert.NotNil(t, err)
}

func TestRemoveStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	createdEvents, err := aggregate.RemoveStorage(storageID, "test")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageRemoved)
	assert.True(t, ok)
	assert.Equal(t, v.Version, 1)
	assert.Equal(t, v.Reason, "test")
}

func TestRemoveStorageForRemovedStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate(true)
	_, err := aggregate.RemoveStorage(storageID, "")
	assert.NotNil(t, err)
}

func TestRemoveStorageNoReason(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	_, err := aggregate.RemoveStorage(storageID, "")
	assert.NotNil(t, err)
}

func TestRenameStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	createdEvents, err := aggregate.SetStorageName(storageID, "storage name set", "test")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageNameSet)
	assert.True(t, ok)
	assert.Equal(t, v.Name, "storage name set")
	assert.Equal(t, v.Reason, "test")
}

func TestRenameStorageForRemovedStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate(true)
	_, err := aggregate.SetStorageName(storageID, "Test storage name set", "")
	assert.NotNil(t, err)
}

func TestRenameStorageWithoutName(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	_, err := aggregate.SetStorageName(storageID, "", "")
	assert.NotNil(t, err)
}

func TestRelocatedStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	createdEvents, err := aggregate.SetStorageLocation(storageID, "location set", "test")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageLocationSet)
	assert.True(t, ok)
	assert.Equal(t, v.Location, "location set")
	assert.Equal(t, v.Reason, "test")
}

func TestRelocatedStorageForRemovedStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate(true)
	_, err := aggregate.SetStorageLocation(storageID, "location set", "")
	assert.NotNil(t, err)
}

func TestRelocateStorageWithoutLocation(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	_, err := aggregate.SetStorageLocation(storageID, "", "")
	assert.NotNil(t, err)
}

func TestOnStorageCreated(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	storageID := uuid.New()
	storageCreated := events.NewStorageCreated(uuid.New(), 1, storageID)
	err := aggregate.On(storageCreated)
	assert.Nil(t, err)
	assert.Len(t, aggregate.Storages, 1)
	assert.Equal(t, aggregate.AggregateVersion(), 1)
	storage := aggregate.GetStorageByID(storageID)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.ID, storageID)
	assert.NotZero(t, storage.CreatedAt)
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnStorageRemoved(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	storageRemoved := events.NewStorageRemoved(uuid.New(), 3, storageID, "test")
	err := aggregate.On(storageRemoved)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 3)
	storage := aggregate.GetStorageByID(storageID)
	assert.NotNil(t, storage)
	assert.True(t, storage.Removed)
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnStorageRenamed(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	storageNameSet := events.NewStorageNameSet(uuid.New(), 5, storageID, "storage set", "test")
	err := aggregate.On(storageNameSet)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 5)
	storage := aggregate.GetStorageByID(storageID)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.Name, "storage set")
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnStorageLocationSet(t *testing.T) {
	storageID, aggregate := newTestAggregate(false)
	storageLocationSet := events.NewStorageLocationSet(uuid.New(), 7, storageID, "location set", "test")
	err := aggregate.On(storageLocationSet)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 7)
	storage := aggregate.GetStorageByID(storageID)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.Location, "location set")
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnUnknownEvent(t *testing.T) {
	_, aggregate := newTestAggregate(false)
	unknownEvent := &UnknownEvent{EventModel: common.EventModel{ID: uuid.New(), Version: 4, At: time.Now()}}
	err := aggregate.On(unknownEvent)
	assert.NotNil(t, err)
}
