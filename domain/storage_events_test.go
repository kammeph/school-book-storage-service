package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
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
			eventType:       domain.StorageAdded,
		},
		{
			name:            "storage removed",
			storageID:       uuid.NewString(),
			storageName:     "storage",
			storageLocation: "location",
			eventType:       domain.StorageRemoved,
		},
		{
			name:        "storage renamed",
			storageID:   uuid.NewString(),
			storageName: "storage",
			reason:      "test",
			eventType:   domain.StorageRenamed,
		},
		{
			name:            "storage relocated",
			storageID:       uuid.NewString(),
			storageLocation: "location",
			reason:          "test",
			eventType:       domain.StorageRelocated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := domain.NewSchoolStorageAggregate()
			var event domain.Event
			var err error
			switch test.eventType {
			case domain.StorageAdded:
				event, err = domain.NewStorageAdded(
					aggregate,
					test.storageID,
					test.storageName,
					test.storageLocation,
				)
			case domain.StorageRemoved:
				event, err = domain.NewStorageRemoved(
					aggregate,
					test.storageID,
					test.reason,
				)
			case domain.StorageRenamed:
				event, err = domain.NewStorageRenamed(
					aggregate,
					test.storageID,
					test.storageName,
					test.reason,
				)
			case domain.StorageRelocated:
				event, err = domain.NewStorageRelocated(
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
