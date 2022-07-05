package storage_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/stretchr/testify/assert"
)

type UnknownEvent struct {
	common.EventModel
}

func newTestAggregate() (uuid.UUID, storage.SchoolAggregateRoot) {
	schoolID := uuid.New()
	storageID := uuid.New()
	aggregate := storage.NewSchoolAggregateRootWithID(schoolID)
	aggregate.School.Storages = append(
		aggregate.School.Storages,
		storage.Storage{
			ID:       storageID,
			Name:     "storage",
			Location: "location"})
	return storageID, aggregate
}

func TestAddStorage(t *testing.T) {
	aggregate := storage.NewSchoolAggregateRoot()
	err := aggregate.AddStorage("storage", "location")
	createdEvents := aggregate.DomainEvents()
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 3)
	created, createdOk := createdEvents[0].(storage.StorageCreated)
	assert.True(t, createdOk)
	assert.Equal(t, created.Version, 1)
	nameSet, nameSetOk := createdEvents[1].(storage.StorageNameSet)
	assert.True(t, nameSetOk)
	assert.Equal(t, nameSet.Version, 2)
	assert.Equal(t, nameSet.Name, "storage")
	locationSet, locationOk := createdEvents[2].(storage.StorageLocationSet)
	assert.True(t, locationOk)
	assert.Equal(t, locationSet.Version, 3)
	assert.Equal(t, locationSet.Location, "location")
	assert.Len(t, aggregate.School.Storages, 1)
}

func TestAddStoragForExistingName(t *testing.T) {
	_, aggregate := newTestAggregate()
	err := aggregate.AddStorage("storage", "location")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageAlreadyExistsError("storage", "location"))
}

func TestAddStorageWithoutName(t *testing.T) {
	aggregate := storage.NewSchoolAggregateRoot()
	err := aggregate.AddStorage("", "location")
	assert.NotNil(t, err)
}

func TestAddStorageWithoutLocation(t *testing.T) {
	aggregate := storage.NewSchoolAggregateRoot()
	err := aggregate.AddStorage("storage", "")
	assert.NotNil(t, err)
}

func TestRemoveStorage(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.RemoveStorage(storageID, "test")
	createdEvents := aggregate.DomainEvents()
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(storage.StorageRemoved)
	assert.True(t, ok)
	assert.Equal(t, v.Version, 1)
	assert.Equal(t, v.Reason, "test")
	assert.Len(t, aggregate.School.Storages, 0)
}

func TestRemoveStorageNotExistsError(t *testing.T) {
	_, aggregate := newTestAggregate()
	unknownID := uuid.New()
	err := aggregate.RemoveStorage(unknownID, "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageIDNotFoundError(unknownID))
}

func TestRemoveStorageNoReasonError(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.RemoveStorage(storageID, "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.ReasonNotSpecifiedError)
}

func TestSetStorageName(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageName(storageID, "storage name set", "test")
	createdEvents := aggregate.DomainEvents()
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(storage.StorageNameSet)
	assert.True(t, ok)
	assert.Equal(t, v.Name, "storage name set")
	assert.Equal(t, v.Reason, "test")
	assert.Equal(t, aggregate.School.Storages[0].Name, "storage name set")
}

func TestSetStorageNameStorageAlreadyExistsError(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageName(storageID, "storage", "test")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageAlreadyExistsError("storage", "location"))
}

func TestSetStorageNameNotStorageExistsError(t *testing.T) {
	_, aggregate := newTestAggregate()
	unknownID := uuid.New()
	err := aggregate.SetStorageName(unknownID, "storage name set", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageIDNotFoundError(unknownID))
}

func TestSetStorageNameValueNotSetError(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageName(storageID, "", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageNameNotSetError)
}

func TestSetStorageNameReasonNotSetError(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageLocation(storageID, "storage", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.ReasonNotSpecifiedError)
}

func TestSetStorageLocation(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageLocation(storageID, "location set", "test")
	createdEvents := aggregate.DomainEvents()
	assert.Nil(t, err)
	assert.Len(t, createdEvents, 1)
	v, ok := createdEvents[0].(storage.StorageLocationSet)
	assert.True(t, ok)
	assert.Equal(t, v.Location, "location set")
	assert.Equal(t, v.Reason, "test")
	assert.Equal(t, aggregate.School.Storages[0].Location, "location set")
}

func TestSetStorageLocationStorageNotFoundError(t *testing.T) {
	_, aggregate := newTestAggregate()
	unknownID := uuid.New()
	err := aggregate.SetStorageLocation(unknownID, "location set", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageIDNotFoundError(unknownID))
}

func TestSetStorageLocationValueNotSetError(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageLocation(storageID, "", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.StorageLocationNotSetError)
}

func TestSetStorageLocationReasonNotSetError(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	err := aggregate.SetStorageLocation(storageID, "location set", "")
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.ReasonNotSpecifiedError)
}

func TestOnStorageCreated(t *testing.T) {
	aggregate := storage.SchoolAggregateRoot{}
	storageID := uuid.New()
	storageCreated := storage.NewStorageCreated(aggregate, storageID)
	err := aggregate.On(storageCreated)
	assert.Nil(t, err)
	assert.Len(t, aggregate.School.Storages, 1)
	assert.Equal(t, aggregate.AggregateVersion(), 1)
	assert.NotZero(t, aggregate.School.UpdatedAt)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.Nil(t, err)
	assert.Greater(t, idx, -1)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.ID, storageID)
	assert.NotZero(t, storage.CreatedAt)
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnStorageRemoved(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	storageRemoved := storage.NewStorageRemoved(aggregate, storageID, "test")
	err := aggregate.On(storageRemoved)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 1)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.NotZero(t, aggregate.School.UpdatedAt)
	assert.NotNil(t, err)
	assert.Equal(t, idx, -1)
	assert.Nil(t, storage)
	assert.Len(t, aggregate.School.Storages, 0)
}

func TestOnStorageRenamed(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	storageNameSet := storage.NewStorageNameSet(aggregate, storageID, "storage set", "test")
	err := aggregate.On(storageNameSet)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 1)
	assert.NotZero(t, aggregate.School.UpdatedAt)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.Nil(t, err)
	assert.Greater(t, idx, -1)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.Name, "storage set")
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnStorageLocationSet(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	storageLocationSet := storage.NewStorageLocationSet(aggregate, storageID, "location set", "test")
	err := aggregate.On(storageLocationSet)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 1)
	assert.NotZero(t, aggregate.School.UpdatedAt)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.Nil(t, err)
	assert.Greater(t, idx, -1)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.Location, "location set")
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnUnknownEvent(t *testing.T) {
	_, aggregate := newTestAggregate()
	unknownEvent := &UnknownEvent{EventModel: common.EventModel{ID: uuid.New(), Version: 4, At: time.Now()}}
	err := aggregate.On(unknownEvent)
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.UnknownEventError(unknownEvent))
}
