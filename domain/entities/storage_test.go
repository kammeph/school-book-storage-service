package entities_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/entities"
	"github.com/kammeph/school-book-storage-service/domain/events"
)

func TestAddStorage(t *testing.T) {
	storage := entities.Storage{}
	createdEvents, err := storage.AddStorage("Test storage", "Test location")
	if err != nil {
		t.Errorf("Error while adding storage: %v", err)
	}
	if len(createdEvents) != 1 {
		t.Errorf("Unexpected length of created events %v", len(createdEvents))
	}
	switch storageAdded := createdEvents[0].(type) {
	case events.StorageAdded:
		if storageAdded.EventVersion() != 1 {
			t.Errorf("Version incorrect: %v", storageAdded.EventVersion())
		}
		if storageAdded.Name != "Test storage" {
			t.Errorf("Name not set: %v", storageAdded.Name)
		}
		if storageAdded.Location != "Test location" {
			t.Errorf("Location not set: %v", storageAdded.Location)
		}
		break
	default:
		t.Errorf("Wrong event type: %v", reflect.TypeOf(storageAdded))
	}
}

func TestAddStorageWithoutName(t *testing.T) {
	storage := entities.Storage{}
	_, err := storage.AddStorage("", "Test location")
	if err == nil {
		t.Error("Name not validated")
	}
}

func TestAddStorageWithoutLocation(t *testing.T) {
	storage := entities.Storage{}
	_, err := storage.AddStorage("Test storage", "")
	if err == nil {
		t.Error("Location not validated")
	}
}

func TestRemoveStorage(t *testing.T) {
	storage := entities.Storage{}
	createdEvents, err := storage.RemoveStorage()
	if err != nil {
		t.Errorf("Error while adding storage: %v", err)
	}
	if len(createdEvents) != 1 {
		t.Errorf("Unexpected length of created events %v", len(createdEvents))
	}
	switch storageRemoved := createdEvents[0].(type) {
	case events.StorageRemoved:
		if storageRemoved.EventVersion() != 1 {
			t.Errorf("Version incorrect: %v", storageRemoved.EventVersion())
		}
		break
	default:
		t.Errorf("Wrong event type: %v", reflect.TypeOf(storageRemoved))
	}
}

func TestRemoveStorageForRemovedStorage(t *testing.T) {
	storage := entities.Storage{Removed: true}
	_, err := storage.RemoveStorage()
	if err == nil {
		t.Error("Incorrect validation for removed storage")
	}
}

func TestRenameStorage(t *testing.T) {
	storage := entities.Storage{Name: "Test storage"}
	createdEvents, err := storage.RenameStorage("Test storage renamed")
	if err != nil {
		t.Errorf("Error while adding storage: %v", err)
	}
	if len(createdEvents) != 1 {
		t.Errorf("Unexpected length of created events %v", len(createdEvents))
	}
	switch storageRenamed := createdEvents[0].(type) {
	case events.StorageRenamed:
		if storageRenamed.EventVersion() != 1 {
			t.Errorf("Version incorrect: %v", storageRenamed.EventVersion())
		}
		if storageRenamed.Name == storage.Name {
			t.Errorf("Storage name and event name are equal: %v", storageRenamed.Name)
		}
		if storageRenamed.Name != "Test storage renamed" {
			t.Errorf("Event name incorrect: %v", storageRenamed.Name)
		}
		break
	default:
		t.Errorf("Wrong event type: %v", reflect.TypeOf(storageRenamed))
	}
}

func TestRenameStorageForRemovedStorage(t *testing.T) {
	storage := entities.Storage{Removed: true}
	_, err := storage.RenameStorage("Test storage renamed")
	if err == nil {
		t.Error("Incorrect validation for removed storage")
	}
}

func TestRenameStorageWithoutName(t *testing.T) {
	storage := entities.Storage{}
	_, err := storage.RenameStorage("")
	if err == nil {
		t.Error("Name not validated")
	}
}

func TestRelocatedStorage(t *testing.T) {
	storage := entities.Storage{Location: "Test location"}
	createdEvents, err := storage.RelocateStorage("Test location relocated")
	if err != nil {
		t.Errorf("Error while adding storage: %v", err)
	}
	if len(createdEvents) != 1 {
		t.Errorf("Unexpected length of created events %v", len(createdEvents))
	}
	switch storageRelocated := createdEvents[0].(type) {
	case events.StorageRelocated:
		if storageRelocated.EventVersion() != 1 {
			t.Errorf("Version incorrect: %v", storageRelocated.EventVersion())
		}
		if storageRelocated.Location == storage.Name {
			t.Errorf("Storage location and event location are equal: %v", storageRelocated.Location)
		}
		if storageRelocated.Location == storage.Name {
			t.Errorf("Event location incorrect: %v", storageRelocated.Location)
		}
		break
	default:
		t.Errorf("Wrong event type: %v", reflect.TypeOf(storageRelocated))
	}
}

func TestRelocatedStorageForRemovedStorage(t *testing.T) {
	storage := entities.Storage{Removed: true}
	_, err := storage.RelocateStorage("Test location relocated")
	if err == nil {
		t.Error("Incorrect validation for removed storage")
	}
}

func TestRelocateStorageWithoutLocation(t *testing.T) {
	storage := entities.Storage{}
	_, err := storage.RelocateStorage("")
	if err == nil {
		t.Error("Location not validated")
	}
}

func TestOnStorageAdded(t *testing.T) {
	storage := entities.Storage{}
	storageAdded := &events.StorageAdded{EventModel: common.EventModel{ID: uuid.New(), Version: 1, At: time.Now()}, Name: "Test storage", Location: "Test location"}
	storage.On(storageAdded)
	if storage.ID != storageAdded.AggregateID() {
		t.Error("ID not set")
	}
	if storage.Version != storageAdded.EventVersion() {
		t.Error("Version not set")
	}
	if storage.CreatedAt != storageAdded.EventAt() {
		t.Error("Created at not set")
	}
	if storage.Name != "Test storage" {
		t.Error("Name not set")
	}
	if storage.Location != "Test location" {
		t.Error("Location not set")
	}
}

func TestOnStorageRemoved(t *testing.T) {
	storage := entities.Storage{}
	storageRemoved := &events.StorageRemoved{EventModel: common.EventModel{ID: uuid.New(), Version: 2, At: time.Now()}}
	storage.On(storageRemoved)
	if storage.Version != storageRemoved.EventVersion() {
		t.Error("Version not set")
	}
	if storage.UpdatedAt != storageRemoved.EventAt() {
		t.Error("Updated at not set")
	}
	if !storage.Removed {
		t.Error("Removed not set")
	}
}

func TestOnStorageRenamed(t *testing.T) {
	storage := entities.Storage{Name: "Test storage"}
	storageRenamed := &events.StorageRenamed{EventModel: common.EventModel{ID: uuid.New(), Version: 3, At: time.Now()}, Name: "Test storage renamed"}
	storage.On(storageRenamed)
	if storage.Version != storageRenamed.EventVersion() {
		t.Error("Version not set")
	}
	if storage.UpdatedAt != storageRenamed.EventAt() {
		t.Error("Updated at not set")
	}
	if storage.Name == "Test storage" {
		t.Error("Name not changed")
	}
	if storage.Name != "Test storage renamed" {
		t.Error("Name not set")
	}
}

func TestOnStorageRelocated(t *testing.T) {
	storage := entities.Storage{Location: "Test location"}
	storageRelocated := &events.StorageRelocated{EventModel: common.EventModel{ID: uuid.New(), Version: 4, At: time.Now()}, Location: "Test location relocated"}
	storage.On(storageRelocated)
	if storage.Version != storageRelocated.EventVersion() {
		t.Error("Version not set")
	}
	if storage.UpdatedAt != storageRelocated.EventAt() {
		t.Error("Updated at not set")
	}
	if storage.Location == "Test location" {
		t.Error("Location not changed")
	}
	if storage.Location != "Test location relocated" {
		t.Error("Location not set")
	}
}
