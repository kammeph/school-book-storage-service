package storage

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type SchoolAggregateRoot struct {
	common.AggregateModel
	School School
}

func NewSchoolAggregateRoot() SchoolAggregateRoot {
	return SchoolAggregateRoot{School: NewSchool()}
}

func NewSchoolAggregateRootWithID(id uuid.UUID) SchoolAggregateRoot {
	return SchoolAggregateRoot{School: NewSchoolWithID(id)}
}

func (a *SchoolAggregateRoot) AddStorage(name, location string) ([]common.Event, error) {
	if name == "" {
		return nil, StorageNameNotSetError
	}
	if location == "" {
		return nil, StorageLocationNotSetError
	}
	for _, storage := range a.School.Storages {
		if storage.Name == name && storage.Location == location {
			return nil, StorageAlreadyExistsError(name, location)
		}
	}
	storageID := uuid.New()
	storageCreated := NewStorageCreated(*a, storageID)
	a.On(storageCreated)
	storageNameSet := NewStorageNameSet(*a, storageID, name, "initial create")
	a.On(storageNameSet)
	storageLocationSet := NewStorageLocationSet(*a, storageID, location, "initial create")
	a.On(storageLocationSet)
	createdEvents := []common.Event{storageCreated, storageNameSet, storageLocationSet}
	return createdEvents, nil
}

func (a *SchoolAggregateRoot) RemoveStorage(storageID uuid.UUID, reason string) ([]common.Event, error) {
	_, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return nil, err
	}
	if reason == "" {
		return nil, ReasonNotSpecifiedError
	}
	event := NewStorageRemoved(*a, storageID, reason)
	a.On(event)
	return []common.Event{event}, nil
}

func (a *SchoolAggregateRoot) SetStorageName(storageID uuid.UUID, name string, reason string) ([]common.Event, error) {
	storage, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, StorageNameNotSetError
	}
	if reason == "" {
		return nil, ReasonNotSpecifiedError
	}
	storagesWithName := a.getStoragesByName(name)
	if len(storagesWithName) > 0 {
		for _, s := range storagesWithName {
			if s.Name == name && s.Location == storage.Location {
				return nil, StorageAlreadyExistsError(s.Name, s.Location)
			}
		}
	}
	event := NewStorageNameSet(*a, storageID, name, reason)
	a.On(event)
	return []common.Event{event}, nil
}

func (a *SchoolAggregateRoot) SetStorageLocation(storageID uuid.UUID, location string, reason string) ([]common.Event, error) {
	storage, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return nil, err
	}
	if location == "" {
		return nil, StorageLocationNotSetError
	}
	if reason == "" {
		return nil, ReasonNotSpecifiedError
	}
	storagesWithLocation := a.getStoragesByLocation(location)
	if len(storagesWithLocation) > 0 {
		for _, s := range storagesWithLocation {
			if s.Name == storage.Name && s.Location == location {
				return nil, StorageAlreadyExistsError(s.Name, s.Location)
			}
		}
	}
	event := NewStorageLocationSet(*a, storageID, location, reason)
	a.On(event)
	return []common.Event{event}, nil
}

func (s *SchoolAggregateRoot) On(event common.Event) error {
	switch evt := event.(type) {
	case StorageCreated:
		return s.onStorageCreated(evt)
	case StorageRemoved:
		return s.onStorageRemoved(evt)
	case StorageNameSet:
		return s.onStorageNameSet(evt)
	case StorageLocationSet:
		return s.onStorageLocationSet(evt)
	default:
		return UnknownEventError(event)
	}
}

func (a *SchoolAggregateRoot) onStorageCreated(event StorageCreated) error {
	storage := NewStorage(event.StorageID, event.EventAt())
	a.Version = event.EventVersion()
	a.School.UpdatedAt = event.EventAt()
	a.School.Storages = append(a.School.Storages, storage)
	return nil
}

func (a *SchoolAggregateRoot) onStorageRemoved(event StorageRemoved) error {
	err := a.RemoveStorageByID(event.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	a.School.UpdatedAt = event.EventAt()
	return nil
}

func (a *SchoolAggregateRoot) onStorageNameSet(event StorageNameSet) error {
	storage, _, err := a.GetStorageByID(event.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	a.School.UpdatedAt = event.EventAt()
	storage.UpdatedAt = event.EventAt()
	storage.Name = event.Name
	return nil
}

func (a *SchoolAggregateRoot) onStorageLocationSet(event StorageLocationSet) error {
	storage, _, err := a.GetStorageByID(event.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	a.School.UpdatedAt = event.EventAt()
	storage.UpdatedAt = event.EventAt()
	storage.Location = event.Location
	return nil
}

func (a *SchoolAggregateRoot) GetStorageByID(id uuid.UUID) (*Storage, int, error) {
	for idx, s := range a.School.Storages {
		if s.ID == id {
			return &a.School.Storages[idx], idx, nil
		}
	}
	return nil, -1, StorageIDNotFoundError(id)
}

func (a *SchoolAggregateRoot) RemoveStorageByID(id uuid.UUID) error {
	_, idx, err := a.GetStorageByID(id)
	if err != nil {
		return err
	}
	a.School.Storages = append(a.School.Storages[:idx], a.School.Storages[idx+1:]...)
	return nil
}

func (a SchoolAggregateRoot) getStoragesByName(name string) []Storage {
	storages := []Storage{}
	for _, storage := range a.School.Storages {
		if storage.Name == name {
			storages = append(storages, storage)
		}
	}
	return storages
}

func (a SchoolAggregateRoot) getStoragesByLocation(location string) []Storage {
	storages := []Storage{}
	for _, storage := range a.School.Storages {
		if storage.Location == location {
			storages = append(storages, storage)
		}
	}
	return storages
}
