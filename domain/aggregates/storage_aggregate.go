package aggregates

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/entities"
	"github.com/kammeph/school-book-storage-service/domain/events"
)

type StorageAggregate struct {
	Storage entities.Storage
}

func (s StorageAggregate) AddStorage(name string, location string) ([]common.Event, error) {
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	event := events.StorageAdded{
		EventModel: common.EventModel{ID: uuid.New(), Version: s.Storage.Version + 1, At: time.Now()},
		Name:       name,
		Location:   location}
	return []common.Event{event}, nil
}

func (s StorageAggregate) RemoveStorage() ([]common.Event, error) {
	if s.Storage.Removed {
		return nil, errors.New("Storage was already removed")
	}
	event := events.StorageRemoved{EventModel: common.EventModel{ID: s.Storage.ID, Version: s.Storage.Version + 1, At: time.Now()}}
	return []common.Event{event}, nil
}

func (s StorageAggregate) RenameStorage(name string) ([]common.Event, error) {
	if s.Storage.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	event := events.StorageRenamed{EventModel: common.EventModel{ID: s.Storage.ID, Version: s.Storage.Version + 1, At: time.Now()}, Name: name}
	return []common.Event{event}, nil
}

func (s StorageAggregate) RelocateStorage(location string) ([]common.Event, error) {
	if s.Storage.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	event := events.StorageRelocated{EventModel: common.EventModel{ID: s.Storage.ID, Version: s.Storage.Version + 1, At: time.Now()}, Location: location}
	return []common.Event{event}, nil
}

func (s *StorageAggregate) On(event common.Event) error {
	switch evt := event.(type) {
	case *events.StorageAdded:
		s.onStorageAdded(*evt)
	case *events.StorageRemoved:
		s.onStorageRemoved(*evt)
	case *events.StorageRenamed:
		s.onStorageRenamed(*evt)
	case *events.StorageRelocated:
		s.onStorageRelocated(*evt)
	default:
		return fmt.Errorf("Unhandled event %v", evt)
	}
	return nil
}

func (s *StorageAggregate) onStorageAdded(event events.StorageAdded) {
	s.Storage.ID = event.AggregateID()
	s.Storage.Version = event.EventVersion()
	s.Storage.CreatedAt = event.EventAt()
	s.Storage.Name = event.Name
	s.Storage.Location = event.Location
}

func (s *StorageAggregate) onStorageRemoved(event events.StorageRemoved) {
	s.Storage.Version = event.EventVersion()
	s.Storage.UpdatedAt = event.EventAt()
	s.Storage.Removed = true
}

func (s *StorageAggregate) onStorageRenamed(event events.StorageRenamed) {
	s.Storage.Version = event.EventVersion()
	s.Storage.UpdatedAt = event.EventAt()
	s.Storage.Name = event.Name
}

func (s *StorageAggregate) onStorageRelocated(event events.StorageRelocated) {
	s.Storage.Version = event.EventVersion()
	s.Storage.UpdatedAt = event.EventAt()
	s.Storage.Location = event.Location
}
