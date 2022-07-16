package storage_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
	"github.com/stretchr/testify/assert"
)

func TestEventCreation(t *testing.T) {
	tests := []struct {
		name            string
		storageID       string
		storageName     string
		storageLocation string
		reason          string
		eventType       string
	}{
		{
			name:            "storage added",
			storageID:       uuid.NewString(),
			storageName:     "storage",
			storageLocation: "location",
			eventType:       storage.StorageAdded,
		},
		{
			name:            "storage removed",
			storageID:       uuid.NewString(),
			storageName:     "storage",
			storageLocation: "location",
			eventType:       storage.StorageRemoved,
		},
		{
			name:        "storage renamed",
			storageID:   uuid.NewString(),
			storageName: "storage",
			reason:      "test",
			eventType:   storage.StorageRenamed,
		},
		{
			name:            "storage relocated",
			storageID:       uuid.NewString(),
			storageLocation: "location",
			reason:          "test",
			eventType:       storage.StorageRelocated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := storage.NewStorageAggregate()
			var event common.Event
			var err error
			switch test.eventType {
			case storage.StorageAdded:
				event, err = storage.NewStorageAdded(
					aggregate,
					test.storageID,
					test.storageName,
					test.storageLocation,
				)
				break
			case storage.StorageRemoved:
				event, err = storage.NewStorageRemoved(
					aggregate,
					test.storageID,
					test.reason,
				)
				break
			case storage.StorageRenamed:
				event, err = storage.NewStorageRenamed(
					aggregate,
					test.storageID,
					test.storageName,
					test.reason,
				)
				break
			case storage.StorageRelocated:
				event, err = storage.NewStorageRelocated(
					aggregate,
					test.storageID,
					test.storageLocation,
					test.reason,
				)
				break
			default:
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, event)
			assert.Equal(t, aggregate.AggregateID(), event.AggregateID())
			assert.Equal(t, aggregate.AggregateVersion()+1, event.EventVersion())
			assert.NotZero(t, event.EventAt())
			assert.Equal(t, test.eventType, event.EventType())
			assert.NotEqual(t, "", event.EventData())
		})
	}
}
