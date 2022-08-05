package storagedomain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
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
			eventType:       storagedomain.StorageAdded,
		},
		{
			name:            "storage removed",
			storageID:       uuid.NewString(),
			storageName:     "storage",
			storageLocation: "location",
			eventType:       storagedomain.StorageRemoved,
		},
		{
			name:        "storage renamed",
			storageID:   uuid.NewString(),
			storageName: "storage",
			reason:      "test",
			eventType:   storagedomain.StorageRenamed,
		},
		{
			name:            "storage relocated",
			storageID:       uuid.NewString(),
			storageLocation: "location",
			reason:          "test",
			eventType:       storagedomain.StorageRelocated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := storagedomain.NewSchoolStorageAggregate()
			var event domain.Event
			var err error
			switch test.eventType {
			case storagedomain.StorageAdded:
				event, err = storagedomain.NewStorageAdded(
					aggregate,
					test.storageID,
					test.storageName,
					test.storageLocation,
				)
			case storagedomain.StorageRemoved:
				event, err = storagedomain.NewStorageRemoved(
					aggregate,
					test.storageID,
					test.reason,
				)
			case storagedomain.StorageRenamed:
				event, err = storagedomain.NewStorageRenamed(
					aggregate,
					test.storageID,
					test.storageName,
					test.reason,
				)
			case storagedomain.StorageRelocated:
				event, err = storagedomain.NewStorageRelocated(
					aggregate,
					test.storageID,
					test.storageLocation,
					test.reason,
				)
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
