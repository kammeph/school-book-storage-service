package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/stretchr/testify/assert"
)

type UnknownEvent struct {
	domain.EventModel
}

func TestOn(t *testing.T) {
	storageID := uuid.NewString()
	tests := []struct {
		name              string
		eventVersion      int
		eventAt           time.Time
		storageName       string
		storageLocation   string
		reason            string
		eventType         string
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
			eventType:         domain.StorageAdded,
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
			eventType:         domain.StorageAdded,
			err:               domain.ErrStoragesWithIdAlreadyExists(storageID),
			expectError:       false,
			addDefaultStorage: true,
		},
		{
			name:              "on storage removed",
			eventVersion:      7,
			eventAt:           time.Now(),
			storageName:       "storage",
			storageLocation:   "location",
			reason:            "test",
			eventType:         domain.StorageRemoved,
			err:               domain.ErrStorageIDNotFound(storageID),
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
			eventType:         domain.StorageRemoved,
			err:               domain.ErrStorageIDNotFound(storageID),
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
			eventType:         domain.StorageRenamed,
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
			eventType:         domain.StorageRenamed,
			err:               domain.ErrStorageIDNotFound(storageID),
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
			eventType:         domain.StorageRelocated,
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
			eventType:         domain.StorageRelocated,
			err:               domain.ErrStorageIDNotFound(storageID),
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
			eventType:         "UnknownEvent",
			err:               domain.ErrUnknownEvent(&UnknownEvent{}),
			expectError:       true,
			addDefaultStorage: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event domain.Event
			switch test.eventType {
			case domain.StorageAdded:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := domain.StorageAddedEvent{"storageAggregate", storageID, test.storageName, test.storageLocation}
				event.SetJsonData(eventData)
			case domain.StorageRemoved:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := domain.StorageRemovedEvent{storageID, test.storageName}
				event.SetJsonData(eventData)
			case domain.StorageRenamed:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := domain.StorageRenamedEvent{storageID, test.storageName, test.reason}
				event.SetJsonData(eventData)
			case domain.StorageRelocated:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := domain.StorageRelocatedEvent{storageID, test.storageLocation, test.reason}
				event.SetJsonData(eventData)
			default:
				event = &UnknownEvent{}
			}
			aggregate := domain.NewSchoolStorageAggregate()
			if test.addDefaultStorage {
				aggregate.Storages = append(aggregate.Storages, domain.Storage{
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
