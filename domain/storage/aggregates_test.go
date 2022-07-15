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
	storageID := uuid.NewString()
	tests := []struct {
		name              string
		eventVersion      int
		eventAt           time.Time
		storageName       string
		storageLocation   string
		reason            string
		event             common.Event
		err               error
		expectError       bool
		addDefaultStorage bool
		operation         string
	}{
		{
			name:              "on storage added",
			eventVersion:      1,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location",
			reason:            "test",
			event:             storage.StorageAdded{},
			err:               nil,
			expectError:       false,
			addDefaultStorage: false,
			operation:         "add",
		},
		{
			name:              "try add storage twice",
			eventVersion:      1,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location",
			reason:            "test",
			event:             storage.StorageAdded{},
			err:               storage.StoragesWithIdAlreadyExistsError(storageID),
			expectError:       true,
			addDefaultStorage: true,
		},
		{
			name:              "on storage removed",
			eventVersion:      7,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location",
			reason:            "test",
			event:             storage.StorageRemoved{},
			err:               storage.StorageIDNotFoundError(storageID),
			expectError:       false,
			addDefaultStorage: true,
			operation:         "remove",
		},
		{
			name:              "remove non existing storage",
			eventVersion:      34,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location",
			reason:            "test",
			event:             storage.StorageRemoved{},
			err:               storage.StorageIDNotFoundError(storageID),
			expectError:       true,
			addDefaultStorage: false,
			operation:         "remove",
		},
		{
			name:              "on storage renamed",
			eventVersion:      5,
			eventAt:           time.Now(),
			storageName:       "storage renamed",
			storageLocation:   "location",
			reason:            "test",
			event:             storage.StorageRenamed{},
			err:               nil,
			expectError:       false,
			addDefaultStorage: true,
			operation:         "update",
		},
		{
			name:              "rename non existing storage",
			eventVersion:      3,
			eventAt:           time.Now(),
			storageName:       "storage renamed",
			storageLocation:   "location",
			reason:            "test",
			event:             storage.StorageRenamed{},
			err:               storage.StorageIDNotFoundError(storageID),
			expectError:       true,
			addDefaultStorage: false,
		},
		{
			name:              "on storage relocated",
			eventVersion:      40,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location relocated",
			reason:            "test",
			event:             storage.StorageRelocated{},
			err:               nil,
			expectError:       false,
			addDefaultStorage: true,
			operation:         "update",
		},
		{
			name:              "relocate non existing storage",
			eventVersion:      9,
			storageName:       "storage",
			storageLocation:   "location relocated",
			reason:            "test",
			event:             storage.StorageRelocated{},
			err:               storage.StorageIDNotFoundError(storageID),
			expectError:       true,
			addDefaultStorage: false,
		},
		{
			name:              "unknown event",
			eventVersion:      9,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location",
			reason:            "test",
			event:             UnknownEvent{},
			err:               storage.UnknownEventError(UnknownEvent{}),
			expectError:       true,
			addDefaultStorage: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event common.Event
			switch test.event.(type) {
			case storage.StorageAdded:
				event = &storage.StorageAdded{
					EventModel: common.EventModel{
						Version: test.eventVersion,
						At:      test.eventAt,
					},
					StorageID: storageID,
					Name:      test.storageName,
					Location:  test.storageLocation,
				}
				break
			case storage.StorageRemoved:
				event = &storage.StorageRemoved{
					EventModel: common.EventModel{
						Version: test.eventVersion,
						At:      test.eventAt,
					},
					StorageID: storageID,
					Reason:    test.reason,
				}
				break
			case storage.StorageRenamed:
				event = &storage.StorageRenamed{
					EventModel: common.EventModel{
						Version: test.eventVersion,
						At:      test.eventAt,
					},
					StorageID: storageID,
					Name:      test.storageName,
					Reason:    test.reason,
				}
				break
			case storage.StorageRelocated:
				event = &storage.StorageRelocated{
					EventModel: common.EventModel{
						Version: test.eventVersion,
						At:      test.eventAt,
					},
					StorageID: storageID,
					Location:  test.storageLocation,
					Reason:    test.reason,
				}
				break
			default:
				event = test.event
				break
			}
			aggregate := storage.NewStorageAggregateRoot()
			if test.addDefaultStorage {
				aggregate.Storages = append(aggregate.Storages, storage.Storage{
					ID:       storageID,
					Name:     "storage",
					Location: "location",
				})
			}
			err := aggregate.On(event)
			if test.expectError {
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.eventVersion, aggregate.Version)
			storage, idx, err := aggregate.GetStorageByID(storageID)
			if test.operation == "remove" {
				assert.Equal(t, test.err, err)
				assert.Equal(t, -1, idx)
				return
			}
			if test.operation == "add" {
				assert.Equal(t, test.eventAt, storage.CreatedAt)
			}
			if test.operation == "update" {
				assert.Equal(t, test.eventAt, storage.UpdatedAt)
			}
			assert.NoError(t, err)
			assert.Greater(t, idx, -1)
			assert.NotNil(t, storage)
			assert.Equal(t, test.storageName, storage.Name)
			assert.Equal(t, test.storageLocation, storage.Location)
		})
	}
}
