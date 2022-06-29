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

func TestAddStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	createdEvents, err := aggregate.AddStorage("storage", "location")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageAdded)
	assert.True(t, ok)
	assert.Equal(t, v.Version, 1)
	assert.Equal(t, v.Name, "storage")
	assert.Equal(t, v.Location, "location")
}

func TestAddStorageWithoutName(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	_, err := aggregate.AddStorage("", "location")
	assert.NotNil(t, err)
}

func TestAddStorageWithoutLocation(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	_, err := aggregate.AddStorage("storage", "")
	assert.NotNil(t, err)
}

func TestRemoveStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	createdEvents, err := aggregate.RemoveStorage()
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageRemoved)
	assert.True(t, ok)
	assert.Equal(t, v.Version, 1)
}

func TestRemoveStorageForRemovedStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Removed: true}}
	_, err := aggregate.RemoveStorage()
	assert.NotNil(t, err)
}

func TestRenameStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Name: "storage"}}
	createdEvents, err := aggregate.RenameStorage("storage name set")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageRenamed)
	assert.True(t, ok)
	assert.Equal(t, v.Name, "storage name set")
}

func TestRenameStorageForRemovedStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Removed: true}}
	_, err := aggregate.RenameStorage("Test storage name set")
	assert.NotNil(t, err)
}

func TestRenameStorageWithoutName(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	_, err := aggregate.RenameStorage("")
	assert.NotNil(t, err)
}

func TestRelocatedStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Location: "location"}}
	createdEvents, err := aggregate.RelocateStorage("location set")
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(events.StorageRelocated)
	assert.True(t, ok)
	assert.Equal(t, v.Location, "location set")
}

func TestRelocatedStorageForRemovedStorage(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Removed: true}}
	_, err := aggregate.RelocateStorage("location set")
	assert.NotNil(t, err)
}

func TestRelocateStorageWithoutLocation(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	_, err := aggregate.RelocateStorage("")
	assert.NotNil(t, err)
}

func TestOnStorageAdded(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	storageAdded := &events.StorageAdded{EventModel: common.EventModel{ID: uuid.New(), Version: 1, At: time.Now()}, Name: "storage", Location: "location"}
	err := aggregate.On(storageAdded)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Storage.ID, storageAdded.AggregateID())
	assert.Equal(t, aggregate.Storage.Version, storageAdded.EventVersion())
	assert.Equal(t, aggregate.Storage.CreatedAt, storageAdded.EventAt())
	assert.Equal(t, aggregate.Storage.Name, "storage")
	assert.Equal(t, aggregate.Storage.Location, "location")
}

func TestOnStorageRemoved(t *testing.T) {
	aggregate := aggregates.StorageAggregate{}
	storageRemoved := &events.StorageRemoved{EventModel: common.EventModel{ID: uuid.New(), Version: 2, At: time.Now()}}
	err := aggregate.On(storageRemoved)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Storage.Version, storageRemoved.EventVersion())
	assert.Equal(t, aggregate.Storage.UpdatedAt, storageRemoved.EventAt())
	assert.True(t, aggregate.Storage.Removed)
}

func TestOnStorageRenamed(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Name: "storage"}}
	storageRenamed := &events.StorageRenamed{EventModel: common.EventModel{ID: uuid.New(), Version: 3, At: time.Now()}, Name: "storage name set"}
	err := aggregate.On(storageRenamed)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Storage.Version, storageRenamed.EventVersion())
	assert.Equal(t, aggregate.Storage.UpdatedAt, storageRenamed.EventAt())
	assert.Equal(t, aggregate.Storage.Name, "storage name set")
}

func TestOnStorageRelocated(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Location: "location"}}
	storageRelocated := &events.StorageRelocated{EventModel: common.EventModel{ID: uuid.New(), Version: 4, At: time.Now()}, Location: "location set"}
	err := aggregate.On(storageRelocated)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Storage.Version, storageRelocated.EventVersion())
	assert.Equal(t, aggregate.Storage.UpdatedAt, storageRelocated.EventAt())
	assert.Equal(t, aggregate.Storage.Location, "location set")
}

func TestOnUnknownEvent(t *testing.T) {
	aggregate := aggregates.StorageAggregate{Storage: entities.Storage{Location: "location"}}
	unknownEvent := &UnknownEvent{EventModel: common.EventModel{ID: uuid.New(), Version: 4, At: time.Now()}}
	err := aggregate.On(unknownEvent)
	assert.NotNil(t, err)
}
