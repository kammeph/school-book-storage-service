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

func newTestAggregate() (string, storage.StorageAggregateRoot) {
	storageID := uuid.New().String()
	aggregate := storage.NewStorageAggregateRoot()
	aggregate.Storages = append(
		aggregate.Storages,
		storage.Storage{
			ID:       storageID,
			Name:     "storage",
			Location: "location",
		})
	return storageID, aggregate
}

func newTestAggregateWithStorages(storages []storage.Storage) storage.StorageAggregateRoot {
	aggregate := storage.NewStorageAggregateRoot()
	aggregate.Storages = storages
	return aggregate
}

func TestAddStorage(t *testing.T) {
	tests := []struct {
		name            string
		storages        []storage.Storage
		storageName     string
		storageLocation string
		err             error
		expectError     bool
	}{
		{
			name:            "add storage",
			storages:        []storage.Storage{},
			storageName:     "storage",
			storageLocation: "location",
			err:             nil,
			expectError:     false,
		},
		{
			name:            "add storage without name",
			storages:        []storage.Storage{},
			storageName:     "",
			storageLocation: "location",
			err:             storage.StorageNameNotSetError,
			expectError:     true,
		},
		{
			name:            "add storage without location",
			storages:        []storage.Storage{},
			storageName:     "storage",
			storageLocation: "",
			err:             storage.StorageLocationNotSetError,
			expectError:     true,
		},
		{
			name: "storage already exists",
			storages: []storage.Storage{{
				ID:       uuid.NewString(),
				Name:     "storage",
				Location: "location",
			}},
			storageName:     "storage",
			storageLocation: "location",
			err:             storage.StorageAlreadyExistsError("storage", "location"),
			expectError:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := newTestAggregateWithStorages(test.storages)
			storageID, err := aggregate.AddStorage(test.storageName, test.storageLocation)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			v, ok := event.(*storage.StorageAdded)
			assert.True(t, ok)
			assert.Equal(t, storageID, v.StorageID)
			assert.Equal(t, test.storageName, v.Name)
			assert.Equal(t, test.storageLocation, v.Location)
		})
	}
}

func TestRemoveStorage(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name        string
		storages    []storage.Storage
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "remove storage",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name:        "error when removing not existing storage",
			storages:    []storage.Storage{},
			reason:      "test",
			err:         storage.StorageIDNotFoundError(storageID),
			expectError: true,
		},
		{
			name: "error when removing without a reason",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			reason:      "",
			err:         storage.ReasonNotSpecifiedError,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := newTestAggregateWithStorages(test.storages)
			err := aggregate.RemoveStorage(storageID, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			v, ok := event.(*storage.StorageRemoved)
			assert.True(t, ok)
			assert.Equal(t, v.Version, 1)
			assert.Equal(t, v.Reason, test.reason)
			assert.Len(t, aggregate.Storages, 0)
		})
	}
}

func TestRenameStorage(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name        string
		storages    []storage.Storage
		storageName string
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "rename storage",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "renamed",
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name: "storage with same name and location exists",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "storage",
			reason:      "test",
			err:         storage.StorageAlreadyExistsError("storage", "location"),
			expectError: true,
		},
		{
			name: "storage name not set",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "",
			reason:      "test",
			err:         storage.StorageNameNotSetError,
			expectError: true,
		},
		{
			name: "reason not specified",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageName: "renamed",
			reason:      "",
			err:         storage.ReasonNotSpecifiedError,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := newTestAggregateWithStorages(test.storages)
			err := aggregate.RenameStorage(storageID, test.storageName, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			v, ok := event.(*storage.StorageRenamed)
			assert.True(t, ok)
			assert.Equal(t, test.storageName, v.Name)
			assert.Equal(t, test.reason, v.Reason)
			assert.Equal(t, test.storageName, aggregate.Storages[0].Name)
		})
	}
}

func TestSetStorageLocation(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name            string
		storages        []storage.Storage
		storageLocation string
		reason          string
		err             error
		expectError     bool
	}{
		{
			name: "rename storage",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "relocated",
			reason:          "test",
			err:             nil,
			expectError:     false,
		},
		{
			name: "storage with same name and location exists",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "location",
			reason:          "test",
			err:             storage.StorageAlreadyExistsError("storage", "location"),
			expectError:     true,
		},
		{
			name: "storage name not set",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "",
			reason:          "test",
			err:             storage.StorageLocationNotSetError,
			expectError:     true,
		},
		{
			name: "reason not specified",
			storages: []storage.Storage{{
				ID:       storageID,
				Name:     "storage",
				Location: "location",
			}},
			storageLocation: "relocated",
			reason:          "",
			err:             storage.ReasonNotSpecifiedError,
			expectError:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := newTestAggregateWithStorages(test.storages)
			err := aggregate.RelocateStorage(storageID, test.storageLocation, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			v, ok := event.(*storage.StorageRelocated)
			assert.True(t, ok)
			assert.Equal(t, test.storageLocation, v.Location)
			assert.Equal(t, test.reason, v.Reason)
			assert.Equal(t, test.storageLocation, aggregate.Storages[0].Location)
		})
	}
}

func TestOnStorageCreated(t *testing.T) {
	aggregate := storage.StorageAggregateRoot{}
	storageID := uuid.NewString()
	name, location := "storage", "location"
	storageCreated := storage.NewStorageAdded(aggregate, storageID, name, location)
	err := aggregate.On(storageCreated)
	assert.Nil(t, err)
	assert.Len(t, aggregate.Storages, 1)
	assert.Equal(t, aggregate.AggregateVersion(), 1)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.Nil(t, err)
	assert.Greater(t, idx, -1)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.ID, storageID)
	assert.Equal(t, storage.Name, name)
	assert.Equal(t, storage.Location, location)
	assert.NotZero(t, storage.CreatedAt)
	assert.Zero(t, storage.UpdatedAt)
}

func TestOnStorageRemoved(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	storageRemoved := storage.NewStorageRemoved(aggregate, storageID, "test")
	err := aggregate.On(storageRemoved)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 1)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.NotNil(t, err)
	assert.Equal(t, idx, -1)
	assert.Nil(t, storage)
	assert.Len(t, aggregate.Storages, 0)
}

func TestOnStorageRenamed(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	storageNameSet := storage.NewStorageRenamed(aggregate, storageID, "storage set", "test")
	err := aggregate.On(storageNameSet)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 1)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.Nil(t, err)
	assert.Greater(t, idx, -1)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.Name, "storage set")
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnStorageRelocated(t *testing.T) {
	storageID, aggregate := newTestAggregate()
	storageLocationSet := storage.NewStorageRelocated(aggregate, storageID, "location set", "test")
	err := aggregate.On(storageLocationSet)
	assert.Nil(t, err)
	assert.Equal(t, aggregate.Version, 1)
	storage, idx, err := aggregate.GetStorageByID(storageID)
	assert.Nil(t, err)
	assert.Greater(t, idx, -1)
	assert.NotNil(t, storage)
	assert.Equal(t, storage.Location, "location set")
	assert.NotZero(t, storage.UpdatedAt)
}

func TestOnUnknownEvent(t *testing.T) {
	_, aggregate := newTestAggregate()
	unknownEvent := &UnknownEvent{EventModel: common.EventModel{ID: uuid.New().String(), Version: 4, At: time.Now()}}
	err := aggregate.On(unknownEvent)
	assert.NotNil(t, err)
	assert.Equal(t, err, storage.UnknownEventError(unknownEvent))
}
