package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/events"
)

type Storage struct {
	ID        uuid.UUID
	Version   int
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Removed   bool
}

func (s Storage) AddStorage(name string, location string) ([]common.Event, error) {
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	event := events.StorageAdded{EventModel: common.EventModel{ID: uuid.New(), Version: s.Version + 1, At: time.Now()}, Name: name, Location: location}
	return []common.Event{event}, nil
}

func (s Storage) RemoveStorage() ([]common.Event, error) {
	if s.Removed {
		return nil, errors.New("Storage was already removed")
	}
	event := events.StorageRemoved{EventModel: common.EventModel{ID: s.ID, Version: s.Version + 1, At: time.Now()}}
	return []common.Event{event}, nil
}

func (s Storage) RenameStorage(name string) ([]common.Event, error) {
	if s.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	event := events.StorageRenamed{EventModel: common.EventModel{ID: s.ID, Version: s.Version + 1, At: time.Now()}, Name: name}
	return []common.Event{event}, nil
}

func (s Storage) RelocateStorage(location string) ([]common.Event, error) {
	if s.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	event := events.StorageRelocated{EventModel: common.EventModel{ID: s.ID, Version: s.Version + 1, At: time.Now()}, Location: location}
	return []common.Event{event}, nil
}

func (s *Storage) On(event common.Event) error {
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

func (s *Storage) onStorageAdded(event events.StorageAdded) {
	s.ID = event.AggregateID()
	s.Version = event.EventVersion()
	s.CreatedAt = event.EventAt()
	s.Name = event.Name
	s.Location = event.Location
}

func (s *Storage) onStorageRemoved(event events.StorageRemoved) {
	s.Version = event.EventVersion()
	s.UpdatedAt = event.EventAt()
	s.Removed = true
}

func (s *Storage) onStorageRenamed(event events.StorageRenamed) {
	s.Version = event.EventVersion()
	s.UpdatedAt = event.EventAt()
	s.Name = event.Name
}

func (s *Storage) onStorageRelocated(event events.StorageRelocated) {
	s.Version = event.EventVersion()
	s.UpdatedAt = event.EventAt()
	s.Location = event.Location
}
