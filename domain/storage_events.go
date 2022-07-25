package domain

var (
	StorageAdded     = "STORAGE_ADDED"
	StorageRemoved   = "STORAGE_REMOVED"
	StorageRenamed   = "STORAGE_RENAMED"
	StorageRelocated = "STORAGE_RELOCATED"
)

type StorageAddedEvent struct {
	SchoolID  string `json:"schoolId"`
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Location  string `json:"location"`
}

func NewStorageAdded(aggregate *SchoolStorageAggregate, storageID, name, location string) (Event, error) {
	eventData := StorageAddedEvent{
		SchoolID:  aggregate.AggregateID(),
		StorageID: storageID,
		Name:      name,
		Location:  location,
	}
	event := NewEvent(aggregate, StorageAdded)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type StorageRemovedEvent struct {
	StorageID string `json:"storageId"`
	Reason    string `json:"reason"`
}

func NewStorageRemoved(aggregate *SchoolStorageAggregate, storageID string, reason string) (Event, error) {
	eventData := StorageRemovedEvent{
		StorageID: storageID,
		Reason:    reason,
	}
	event := NewEvent(aggregate, StorageRemoved)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type StorageRenamedEvent struct {
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
}

func NewStorageRenamed(aggregate *SchoolStorageAggregate, storageID string, name, reason string) (Event, error) {
	eventData := StorageRenamedEvent{
		StorageID: storageID,
		Name:      name,
		Reason:    reason,
	}
	event := NewEvent(aggregate, StorageRenamed)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type StorageRelocatedEvent struct {
	StorageID string `json:"storageId"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
}

func NewStorageRelocated(aggregate *SchoolStorageAggregate, storageID string, location, reason string) (Event, error) {
	eventData := StorageRelocatedEvent{
		StorageID: storageID,
		Location:  location,
		Reason:    reason,
	}
	event := NewEvent(aggregate, StorageRelocated)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}
