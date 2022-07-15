package storage

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type StorageAggregateRoot struct {
	common.AggregateModel
	Storages []Storage
}

func NewStorageAggregateRoot() StorageAggregateRoot {
	return StorageAggregateRoot{Storages: []Storage{}}
}

func (a *StorageAggregateRoot) AddStorage(name, location string) (string, error) {
	if name == "" {
		return "", StorageNameNotSetError
	}
	if location == "" {
		return "", StorageLocationNotSetError
	}
	for _, storage := range a.Storages {
		if storage.Name == name && storage.Location == location {
			return "", StorageAlreadyExistsError(name, location)
		}
	}
	storageID := uuid.NewString()
	event := NewStorageAdded(*a, storageID, name, location)
	a.On(event)
	a.Events = append(a.Events, event)
	return storageID, nil
}

func (a *StorageAggregateRoot) RemoveStorage(storageID string, reason string) error {
	_, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if reason == "" {
		return ReasonNotSpecifiedError
	}
	event := NewStorageRemoved(*a, storageID, reason)
	a.On(event)
	a.Events = append(a.Events, event)
	return nil
}

func (a *StorageAggregateRoot) RenameStorage(storageID string, name string, reason string) error {
	storage, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if name == "" {
		return StorageNameNotSetError
	}
	if reason == "" {
		return ReasonNotSpecifiedError
	}
	storagesWithName := a.getStoragesByName(name)
	if len(storagesWithName) > 0 {
		for _, s := range storagesWithName {
			if s.Name == name && s.Location == storage.Location {
				return StorageAlreadyExistsError(s.Name, s.Location)
			}
		}
	}
	event := NewStorageRenamed(*a, storageID, name, reason)
	a.On(event)
	a.Events = append(a.Events, event)
	return nil
}

func (a *StorageAggregateRoot) RelocateStorage(storageID string, location string, reason string) error {
	storage, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if location == "" {
		return StorageLocationNotSetError
	}
	if reason == "" {
		return ReasonNotSpecifiedError
	}
	storagesWithLocation := a.getStoragesByLocation(location)
	if len(storagesWithLocation) > 0 {
		for _, s := range storagesWithLocation {
			if s.Name == storage.Name && s.Location == location {
				return StorageAlreadyExistsError(s.Name, s.Location)
			}
		}
	}
	event := NewStorageRelocated(*a, storageID, location, reason)
	a.On(event)
	a.Events = append(a.Events, event)
	return nil
}

func (s *StorageAggregateRoot) On(event common.Event) error {
	switch evt := event.(type) {
	case *StorageAdded:
		return s.onStorageAdded(evt)
	case *StorageRemoved:
		return s.onStorageRemoved(evt)
	case *StorageRenamed:
		return s.onStorageRenamed(evt)
	case *StorageRelocated:
		return s.onStorageRelocated(evt)
	default:
		return UnknownEventError(event)
	}
}

func (a *StorageAggregateRoot) onStorageAdded(event *StorageAdded) error {
	_, idx, _ := a.GetStorageByID(event.StorageID)
	if idx > -1 {
		return StoragesWithIdAlreadyExistsError(event.StorageID)
	}
	storage := NewStorage(event.StorageID, event.Name, event.Location, event.EventAt())
	a.Version = event.EventVersion()
	a.Storages = append(a.Storages, storage)
	return nil
}

func (a *StorageAggregateRoot) onStorageRemoved(event *StorageRemoved) error {
	err := a.RemoveStorageByID(event.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	return nil
}

func (a *StorageAggregateRoot) onStorageRenamed(event *StorageRenamed) error {
	storage, _, err := a.GetStorageByID(event.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	storage.UpdatedAt = event.EventAt()
	storage.Name = event.Name
	return nil
}

func (a *StorageAggregateRoot) onStorageRelocated(event *StorageRelocated) error {
	storage, _, err := a.GetStorageByID(event.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	storage.UpdatedAt = event.EventAt()
	storage.Location = event.Location
	return nil
}

func (a *StorageAggregateRoot) GetStorageByID(id string) (*Storage, int, error) {
	for idx, s := range a.Storages {
		if s.ID == id {
			return &a.Storages[idx], idx, nil
		}
	}
	return nil, -1, StorageIDNotFoundError(id)
}

func (a *StorageAggregateRoot) RemoveStorageByID(id string) error {
	_, idx, err := a.GetStorageByID(id)
	if err != nil {
		return err
	}
	a.Storages = append(a.Storages[:idx], a.Storages[idx+1:]...)
	return nil
}

func (a StorageAggregateRoot) GetStorageByName(name string) (*Storage, error) {
	storages := a.getStoragesByName(name)
	if len(storages) == 0 {
		return nil, StorageByNameNotFoundError(name)
	}
	if len(storages) > 1 {
		return nil, MultipleStoragesWithNameFoundError(name)
	}
	return &storages[0], nil
}

func (a StorageAggregateRoot) getStoragesByName(name string) []Storage {
	storages := []Storage{}
	for _, storage := range a.Storages {
		if storage.Name == name {
			storages = append(storages, storage)
		}
	}
	return storages
}

func (a StorageAggregateRoot) getStoragesByLocation(location string) []Storage {
	storages := []Storage{}
	for _, storage := range a.Storages {
		if storage.Location == location {
			storages = append(storages, storage)
		}
	}
	return storages
}
