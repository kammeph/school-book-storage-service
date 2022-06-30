package aggregates

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/domain/entities"
	"github.com/kammeph/school-book-storage-service/domain/events"
)

type StorageAggregate struct {
	common.AggregateModel
	Storages []entities.Storage
}

func NewStorageAggregate() StorageAggregate {
	return StorageAggregate{Storages: []entities.Storage{}}
}

func (a StorageAggregate) AddStorage(name string, location string) ([]common.Event, error) {
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	for _, storage := range a.Storages {
		if storage.Name == name && storage.Location == location {
			return nil, fmt.Errorf("Storage with name %s in location %s already exists", name, location)
		}
	}
	nextVersion := a.Version + 1
	storageID := uuid.New()
	storageCreated := events.NewStorageCreated(a.AggregateID(), nextVersion, storageID)
	nextVersion += 1
	storageNameSet := events.NewStorageNameSet(a.AggregateID(), nextVersion, storageID, name, "initial create")
	nextVersion += 1
	storageLocationSet := events.NewStorageLocationSet(a.AggregateID(), nextVersion, storageID, location, "initial create")
	createdEvents := []common.Event{storageCreated, storageNameSet, storageLocationSet}
	return createdEvents, nil
}

func (a StorageAggregate) RemoveStorage(storageID uuid.UUID, reason string) ([]common.Event, error) {
	storage := a.GetStorageByID(storageID)
	if storage == nil {
		return nil, fmt.Errorf("The storage with the ID %s does not exist.", storageID.String())
	}
	if storage.Removed {
		return nil, fmt.Errorf("The storage with the ID %s is already removed.", storageID.String())
	}
	if reason == "" {
		return nil, errors.New("No reason specified")
	}
	event := events.NewStorageRemoved(a.AggregateID(), a.AggregateVersion()+1, storageID, reason)
	return []common.Event{event}, nil
}

func (a StorageAggregate) SetStorageName(storageID uuid.UUID, name string, reason string) ([]common.Event, error) {
	storage := a.GetStorageByID(storageID)
	if storage == nil {
		return nil, fmt.Errorf("The storage with the ID %s does not exist.", storageID.String())
	}
	if storage.Removed {
		return nil, fmt.Errorf("The storage with the ID %s is already removed.", storageID.String())
	}
	if name == "" {
		return nil, errors.New("Storage name unknown")
	}
	if reason == "" {
		return nil, errors.New("No reason specified")
	}
	storagesWithName := a.getStoragesByName(name)
	for _, s := range a.Storages {
		if s.ID == storageID {
			storage = &s
		}
		if s.Name == name {
			return nil, fmt.Errorf("There is already a storage with the name %s", name)
		}
	}
	if len(storagesWithName) > 0 {
		for _, s := range storagesWithName {
			if s.Name == name && s.Location == storage.Location {
				return nil, fmt.Errorf("Storage with name %s in location %s already exists", s.Name, s.Location)
			}
		}
	}
	event := events.NewStorageNameSet(a.AggregateID(), a.AggregateVersion()+1, storageID, name, reason)
	return []common.Event{event}, nil
}

func (a StorageAggregate) SetStorageLocation(storageID uuid.UUID, location string, reason string) ([]common.Event, error) {
	storage := a.GetStorageByID(storageID)
	if storage == nil {
		return nil, fmt.Errorf("The storage with the ID %s does not exist.", storageID.String())
	}
	if storage.Removed {
		return nil, fmt.Errorf("The storage with the ID %s is already removed.", storageID.String())
	}
	if location == "" {
		return nil, errors.New("Storage location unknown")
	}
	if reason == "" {
		return nil, errors.New("No reason specified")
	}
	storagesWithLocation := a.getStoragesByLocation(location)
	if len(storagesWithLocation) > 0 {
		for _, s := range storagesWithLocation {
			if s.Name == storage.Name && s.Location == location {
				return nil, fmt.Errorf("Storage with name %s in location %s already exists", s.Name, s.Location)
			}
		}
	}
	event := events.NewStorageLocationSet(a.AggregateID(), a.AggregateVersion()+1, storageID, location, reason)
	return []common.Event{event}, nil
}

func (s *StorageAggregate) On(event common.Event) error {
	switch evt := event.(type) {
	case events.StorageCreated:
		s.onStorageCreated(evt)
	case events.StorageRemoved:
		s.onStorageRemoved(evt)
	case events.StorageNameSet:
		s.onStorageNameSet(evt)
	case events.StorageLocationSet:
		s.onStorageLocationSet(evt)
	default:
		return fmt.Errorf("Unhandled event %v", evt)
	}
	return nil
}

func (a *StorageAggregate) onStorageCreated(event events.StorageCreated) {
	storage := entities.NewStorage(event.StorageID, event.EventAt())
	a.Storages = append(a.Storages, storage)
	a.Version = event.EventVersion()
}

func (a *StorageAggregate) onStorageRemoved(event events.StorageRemoved) {
	storage := a.GetStorageByID(event.StorageID)
	storage.UpdatedAt = event.EventAt()
	storage.Removed = true
	a.Version = event.EventVersion()
}

func (a *StorageAggregate) onStorageNameSet(event events.StorageNameSet) {
	storage := a.GetStorageByID(event.StorageID)
	storage.UpdatedAt = event.EventAt()
	storage.Name = event.Name
	a.Version = event.EventVersion()
}

func (a *StorageAggregate) onStorageLocationSet(event events.StorageLocationSet) {
	storage := a.GetStorageByID(event.StorageID)
	storage.UpdatedAt = event.EventAt()
	storage.Location = event.Location
	a.Version = event.EventVersion()
}

func (a *StorageAggregate) GetStorageByID(id uuid.UUID) *entities.Storage {
	for idx, s := range a.Storages {
		if s.ID == id {
			return &a.Storages[idx]
		}
	}
	return nil
}

func (a StorageAggregate) getStoragesByName(name string) []entities.Storage {
	storages := []entities.Storage{}
	for _, storage := range a.Storages {
		if storage.Name == name {
			storages = append(storages, storage)
		}
	}
	return storages
}

func (a StorageAggregate) getStoragesByLocation(location string) []entities.Storage {
	storages := []entities.Storage{}
	for _, storage := range a.Storages {
		if storage.Location == location {
			storages = append(storages, storage)
		}
	}
	return storages
}
