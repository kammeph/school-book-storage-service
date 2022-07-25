package storage_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	application "github.com/kammeph/school-book-storage-service/application/storage"
	"github.com/kammeph/school-book-storage-service/domain/common"
	domain "github.com/kammeph/school-book-storage-service/domain/storage"
	infrastructure "github.com/kammeph/school-book-storage-service/infrastructure/storage"
	"github.com/stretchr/testify/assert"
)

var (
	eventHandler = application.NewStorageEventHandler(infrastructure.NewMemoryRepository())
	storageAdded = common.EventModel{
		ID:      "school1",
		Version: 1,
		At:      time.Now(),
		Type:    domain.StorageAdded,
		Data:    "{\"schoolId\":\"school1\",\"storageId\":\"storage1\",\"name\":\"closet 2\",\"location\":\"room 101\"}",
	}
	storageRenamed = common.EventModel{
		ID:      "school1",
		Version: 2,
		At:      time.Now(),
		Type:    domain.StorageRenamed,
		Data:    "{\"storageId\":\"storage1\",\"name\":\"closet renamed\",\"reason\":\"test\"}",
	}
	storageRelocated = common.EventModel{
		ID:      "school1",
		Version: 3,
		At:      time.Now(),
		Type:    domain.StorageRelocated,
		Data:    "{\"storageId\":\"storage1\",\"location\":\"location 2\",\"reason\":\"test\"}",
	}
	storageRemoved = common.EventModel{
		ID:      "school1",
		Version: 4,
		At:      time.Now(),
		Type:    domain.StorageRemoved,
		Data:    "{\"storageId\":\"storage1\",\"reason\":\"test\"}",
	}
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name  string
		event common.Event
	}{
		{
			name:  "storage added",
			event: &storageAdded,
		},
		{
			name:  "storage renamed",
			event: &storageRenamed,
		},
		{
			name:  "storage relocated",
			event: &storageRelocated,
		},
		{
			name:  "storage removed",
			event: &storageRemoved,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eventBytes, _ := json.Marshal(test.event)
			assert.NotPanics(t, func() { eventHandler.Handle(context.Background(), eventBytes) })
		})
	}
}
