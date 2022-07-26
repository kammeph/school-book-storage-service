package storagedomain

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
)

func (a *SchoolStorageAggregate) AddStorage(name, location string) (string, error) {
	if name == "" {
		return "", ErrStorageNameNotSet
	}
	if location == "" {
		return "", ErrStorageLocationNotSet
	}
	for _, storage := range a.Storages {
		if storage.Name == name && storage.Location == location {
			return "", ErrStorageAlreadyExists(name, location)
		}
	}
	storageID := uuid.NewString()
	event, err := NewStorageAdded(a, storageID, name, location)
	if err != nil {
		return "", err
	}
	if err := a.Apply(event); err != nil {
		return "", err
	}
	return storageID, nil
}

func (a *SchoolStorageAggregate) RemoveStorage(storageID string, reason string) error {
	_, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	event, err := NewStorageRemoved(a, storageID, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *SchoolStorageAggregate) RenameStorage(storageID string, name string, reason string) error {
	storage, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if name == "" {
		return ErrStorageNameNotSet
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	storagesWithName := a.getStoragesByName(name)
	if len(storagesWithName) > 0 {
		for _, s := range storagesWithName {
			if s.Name == name && s.Location == storage.Location {
				return ErrStorageAlreadyExists(s.Name, s.Location)
			}
		}
	}
	event, err := NewStorageRenamed(a, storageID, name, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *SchoolStorageAggregate) RelocateStorage(storageID string, location string, reason string) error {
	storage, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if location == "" {
		return ErrStorageLocationNotSet
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	storagesWithLocation := a.getStoragesByLocation(location)
	if len(storagesWithLocation) > 0 {
		for _, s := range storagesWithLocation {
			if s.Name == storage.Name && s.Location == location {
				return ErrStorageAlreadyExists(s.Name, s.Location)
			}
		}
	}
	event, err := NewStorageRelocated(a, storageID, location, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}
