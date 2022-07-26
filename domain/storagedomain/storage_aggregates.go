package storagedomain

import "github.com/kammeph/school-book-storage-service/domain"

type SchoolStorageAggregate struct {
	*domain.AggregateModel
	Storages []Storage
}

func NewSchoolStorageAggregate() *SchoolStorageAggregate {
	aggregate := &SchoolStorageAggregate{
		Storages: []Storage{},
	}
	model := domain.NewAggregateModel(aggregate.On)
	aggregate.AggregateModel = &model
	return aggregate
}

func NewSchoolStorageAggregateWithID(id string) *SchoolStorageAggregate {
	aggregate := NewSchoolStorageAggregate()
	aggregate.ID = id
	return aggregate
}

func (s *SchoolStorageAggregate) On(event domain.Event) error {
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
		return domain.ErrUnknownEvent(event)
	}
}

func (a *SchoolStorageAggregate) onStorageAdded(event domain.Event) error {
	eventData := StorageAddedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	_, idx, _ := a.GetStorageByID(eventData.StorageID)
	if idx > -1 {
		return ErrStoragesWithIdAlreadyExists(eventData.StorageID)
	}
	storage := NewStorage(eventData.StorageID, eventData.Name, eventData.Location, event.EventAt())
	a.Version = event.EventVersion()
	a.Storages = append(a.Storages, storage)
	return nil
}

func (a *SchoolStorageAggregate) onStorageRemoved(event domain.Event) error {
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

func (a *SchoolStorageAggregate) onStorageRenamed(event domain.Event) error {
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

func (a *SchoolStorageAggregate) onStorageRelocated(event domain.Event) error {
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

func (a *SchoolStorageAggregate) GetStorageByID(id string) (*Storage, int, error) {
	for idx, s := range a.Storages {
		if s.ID == id {
			return &a.Storages[idx], idx, nil
		}
	}
	return nil, -1, ErrStorageIDNotFound(id)
}

func (a *SchoolStorageAggregate) RemoveStorageByID(id string) error {
	_, idx, err := a.GetStorageByID(id)
	if err != nil {
		return err
	}
	a.Storages = append(a.Storages[:idx], a.Storages[idx+1:]...)
	return nil
}

func (a SchoolStorageAggregate) GetStorageByName(name string) (*Storage, error) {
	storages := a.getStoragesByName(name)
	if len(storages) == 0 {
		return nil, ErrStorageByNameNotFound(name)
	}
	if len(storages) > 1 {
		return nil, ErrMultipleStoragesWithNameFound(name)
	}
	return &storages[0], nil
}

func (a SchoolStorageAggregate) getStoragesByName(name string) []Storage {
	storages := []Storage{}
	for _, storage := range a.Storages {
		if storage.Name == name {
			storages = append(storages, storage)
		}
	}
	return storages
}

func (a SchoolStorageAggregate) getStoragesByLocation(location string) []Storage {
	storages := []Storage{}
	for _, storage := range a.Storages {
		if storage.Location == location {
			storages = append(storages, storage)
		}
	}
	return storages
}
