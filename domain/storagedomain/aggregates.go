package storagedomain

import (
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

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
	if fp.Some(a.Storages, func(s Storage) bool { return s.ID == eventData.StorageID }) {
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
	a.Storages = fp.Remove(a.Storages, func(s Storage) bool { return s.ID == eventData.StorageID })
	a.Version = event.EventVersion()
	return nil
}

func (a *SchoolStorageAggregate) onStorageRenamed(event domain.Event) error {
	eventData := StorageRenamedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	storage := fp.Find(a.Storages, func(s Storage) bool { return s.ID == eventData.StorageID })
	if storage == nil {
		return ErrStorageIDNotFound(eventData.StorageID)
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
	storage := fp.Find(a.Storages, func(s Storage) bool { return s.ID == eventData.StorageID })
	if storage == nil {
		return ErrStorageIDNotFound(eventData.StorageID)
	}
	a.Version = event.EventVersion()
	storage.UpdatedAt = event.EventAt()
	storage.Location = eventData.Location
	return nil
}
