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

func (s StorageAggregate) AddStorage(id uuid.UUID, name string, location string) ([]common.Event, error) {
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	nextVersion := s.Storage.Version + 1
	storageCreated := events.StorageCreated{
		EventModel: common.EventModel{ID: id, Version: nextVersion, At: time.Now()}}
	createdEvents := []common.Event{storageCreated}
	nextVersion += 1
	storageNameSet := events.StorageNameSet{
		EventModel: common.EventModel{ID: id, Version: nextVersion, At: time.Now()},
		Name:       name,
		Reason:     "initial create"}
	createdEvents = append(createdEvents, storageNameSet)
	nextVersion += 1
	storageLocationSet := events.StorageLocationSet{
		EventModel: common.EventModel{ID: id, Version: nextVersion, At: time.Now()},
		Location:   location,
		Reason:     "initial create"}
	createdEvents = append(createdEvents, storageLocationSet)
	return createdEvents, nil
}

func (s StorageAggregate) RemoveStorage(reason string) ([]common.Event, error) {
	if s.Storage.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if reason == "" {
		return nil, errors.New("No reason specified")
	}
	event := events.StorageRemoved{
		EventModel: common.EventModel{ID: s.Storage.ID, Version: s.Storage.Version + 1, At: time.Now()},
		Reason:     reason}
	return []common.Event{event}, nil
}

func (s StorageAggregate) SetStorageName(name string, reason string) ([]common.Event, error) {
	if s.Storage.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	event := events.StorageNameSet{
		EventModel: common.EventModel{ID: s.Storage.ID, Version: s.Storage.Version + 1, At: time.Now()},
		Name:       name,
		Reason:     reason}
	return []common.Event{event}, nil
}

func (s StorageAggregate) SetStorageLocation(location string, reason string) ([]common.Event, error) {
	if s.Storage.Removed {
		return nil, errors.New("Storage was already removed")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	event := events.StorageLocationSet{
		EventModel: common.EventModel{ID: s.Storage.ID, Version: s.Storage.Version + 1, At: time.Now()},
		Location:   location,
		Reason:     reason}
	return []common.Event{event}, nil
}

func (s *StorageAggregate) On(event common.Event) error {
	switch evt := event.(type) {
	case *events.StorageCreated:
		s.onStorageCreated(*evt)
	case *events.StorageRemoved:
		s.onStorageRemoved(*evt)
	case *events.StorageNameSet:
		s.onStorageNameSet(*evt)
	case *events.StorageLocationSet:
		s.onStorageLocationSet(*evt)
	default:
		return fmt.Errorf("Unhandled event %v", evt)
	}
	return nil
}

func (s *StorageAggregate) onStorageCreated(event events.StorageCreated) {
	s.Storage.ID = event.AggregateID()
	s.Storage.Version = event.EventVersion()
	s.Storage.CreatedAt = event.EventAt()
	s.Storage.UpdatedAt = event.EventAt()
}

func (s *StorageAggregate) onStorageRemoved(event events.StorageRemoved) {
	s.Storage.Version = event.EventVersion()
	s.Storage.UpdatedAt = event.EventAt()
	s.Storage.Removed = true
}

func (s *StorageAggregate) onStorageNameSet(event events.StorageNameSet) {
	s.Storage.Version = event.EventVersion()
	s.Storage.UpdatedAt = event.EventAt()
	s.Storage.Name = event.Name
}

func (s *StorageAggregate) onStorageLocationSet(event events.StorageLocationSet) {
	s.Storage.Version = event.EventVersion()
	s.Storage.UpdatedAt = event.EventAt()
	s.Storage.Location = event.Location
}
