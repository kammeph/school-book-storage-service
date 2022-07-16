package storage

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type StorageAggregate struct {
	*common.AggregateModel
	Storages []Storage
}

func NewStorageAggregate() *StorageAggregate {
	aggregate := &StorageAggregate{
		Storages: []Storage{},
	}
	model := common.NewAggregateModel(aggregate.On)
	aggregate.AggregateModel = &model
	return aggregate
}

func (a *StorageAggregate) AddStorage(name, location string) (string, error) {
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
	event, err := NewStorageAdded(a, storageID, name, location)
	if err != nil {
		return "", err
	}
	if err := a.Apply(event); err != nil {
		return "", err
	}
	return storageID, nil
}

func (a *StorageAggregate) RemoveStorage(storageID string, reason string) error {
	_, _, err := a.GetStorageByID(storageID)
	if err != nil {
		return err
	}
	if reason == "" {
		return ReasonNotSpecifiedError
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

func (a *StorageAggregate) RenameStorage(storageID string, name string, reason string) error {
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
	event, err := NewStorageRenamed(a, storageID, name, reason)
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *StorageAggregate) RelocateStorage(storageID string, location string, reason string) error {
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
	event, err := NewStorageRelocated(a, storageID, location, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (s *StorageAggregate) On(event common.Event) error {
	switch event.EventType() {
	case StorageAdded:
		return s.onStorageAdded(event)
	case StorageRemoved:
		return s.onStorageRemoved(event)
	case StorageRenamed:
		return s.onStorageRenamed(event)
	case StorageRelocated:
		return s.onStorageRelocated(event)
	default:
		return UnknownEventError(event)
	}
}

func (a *StorageAggregate) onStorageAdded(event common.Event) error {
	eventData := StorageAddedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	_, idx, _ := a.GetStorageByID(eventData.StorageID)
	if idx > -1 {
		return StoragesWithIdAlreadyExistsError(eventData.StorageID)
	}
	storage := NewStorage(eventData.StorageID, eventData.Name, eventData.Location, event.EventAt())
	a.Version = event.EventVersion()
	a.Storages = append(a.Storages, storage)
	return nil
}

func (a *StorageAggregate) onStorageRemoved(event common.Event) error {
	eventData := StorageRemovedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	if err := a.RemoveStorageByID(eventData.StorageID); err != nil {
		return err
	}
	a.Version = event.EventVersion()
	return nil
}

func (a *StorageAggregate) onStorageRenamed(event common.Event) error {
	eventData := StorageRenamedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	storage, _, err := a.GetStorageByID(eventData.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	storage.UpdatedAt = event.EventAt()
	storage.Name = eventData.Name
	return nil
}

func (a *StorageAggregate) onStorageRelocated(event common.Event) error {
	eventData := StorageRelocatedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	storage, _, err := a.GetStorageByID(eventData.StorageID)
	if err != nil {
		return err
	}
	a.Version = event.EventVersion()
	storage.UpdatedAt = event.EventAt()
	storage.Location = eventData.Location
	return nil
}

func (a *StorageAggregate) GetStorageByID(id string) (*Storage, int, error) {
	for idx, s := range a.Storages {
		if s.ID == id {
			return &a.Storages[idx], idx, nil
		}
	}
	return nil, -1, StorageIDNotFoundError(id)
}

func (a *StorageAggregate) RemoveStorageByID(id string) error {
	_, idx, err := a.GetStorageByID(id)
	if err != nil {
		return err
	}
	a.Storages = append(a.Storages[:idx], a.Storages[idx+1:]...)
	return nil
}

func (a StorageAggregate) GetStorageByName(name string) (*Storage, error) {
	storages := a.getStoragesByName(name)
	if len(storages) == 0 {
		return nil, StorageByNameNotFoundError(name)
	}
	if len(storages) > 1 {
		return nil, MultipleStoragesWithNameFoundError(name)
	}
	return &storages[0], nil
}

func (a StorageAggregate) getStoragesByName(name string) []Storage {
	storages := []Storage{}
	for _, storage := range a.Storages {
		if storage.Name == name {
			storages = append(storages, storage)
		}
	}
	return storages
}

func (a StorageAggregate) getStoragesByLocation(location string) []Storage {
	storages := []Storage{}
	for _, storage := range a.Storages {
		if storage.Location == location {
			storages = append(storages, storage)
		}
	}
	return storages
}
