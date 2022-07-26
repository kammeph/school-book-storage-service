package domain

var (
	SchoolAdded       = "SCHOOL_ADDED"
	SchoolDeactivated = "SCHOOL_DEACTIVATED"
	SchoolRenamed     = "SCHOOL_RENAMED"
)

type SchoolAddedEvent struct {
	SchoolID string `json:"schoolId"`
	Name     string `json:"name"`
}

func NewSchoolAdded(aggregate *SchoolAggregate, schoolID, name string) (Event, error) {
	eventData := SchoolAddedEvent{
		SchoolID: schoolID,
		Name:     name,
	}
	event := NewEvent(aggregate, SchoolAdded)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type SchoolDeactivatedEvent struct {
	SchoolID string `json:"schoolID"`
	Reason   string `json:"reason"`
}

func NewSchoolDeactivated(aggregate *SchoolAggregate, schoolID, reason string) (Event, error) {
	eventData := SchoolDeactivatedEvent{
		SchoolID: schoolID,
		Reason:   reason,
	}
	event := NewEvent(aggregate, SchoolDeactivated)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type SchoolRenamedEvent struct {
	SchoolID string `json:"schoolId"`
	Name     string `json:"name"`
	Reason   string `json:"reason"`
}

func NewSchoolRenamed(aggregate *SchoolAggregate, schoolID, name, reason string) (Event, error) {
	eventData := SchoolRenamedEvent{
		SchoolID: schoolID,
		Name:     name,
		Reason:   reason,
	}
	event := NewEvent(aggregate, SchoolRenamed)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}
