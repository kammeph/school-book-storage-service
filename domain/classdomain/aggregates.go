package classdomain

import (
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

type SchoolClassAggregate struct {
	*domain.AggregateModel
	Classes []Class
}

func NewSchoolClassAggregate() *SchoolClassAggregate {
	aggregate := &SchoolClassAggregate{
		Classes: []Class{},
	}
	model := domain.NewAggregateModel(aggregate.On)
	aggregate.AggregateModel = &model
	return aggregate
}

func (a *SchoolClassAggregate) On(event domain.Event) error {
	switch event.EventType() {
	case ClassCreated:
		return a.onClassCreated(event)
	case NumberOfPupilsIncreased:
		return a.onNumberOfPupilsIncreased(event)
	case NumberOfPupilsDecreased:
		return a.onNumberOfPupilsDecreased(event)
	default:
		return domain.ErrUnknownEvent(event)
	}
}

func (a *SchoolClassAggregate) onClassCreated(event domain.Event) error {
	eventData := ClassCreatedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	if fp.Some(a.Classes, func(c Class) bool { return c.ID == eventData.ClassID }) {
		return ErrApplyEventClassAlreadyExists(event.EventType(), eventData.ClassID)
	}
	class := NewClass(
		eventData.ClassID,
		eventData.Grade,
		eventData.Letter,
		eventData.NumberOfPupils,
		eventData.DateFrom,
		eventData.DateTo,
		event.EventAt())
	a.Version = event.EventVersion()
	a.Classes = append(a.Classes, class)
	return nil
}

func (a *SchoolClassAggregate) onNumberOfPupilsIncreased(event domain.Event) error {
	eventData := NumberOfPupilsIncreasedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	class := fp.Find(a.Classes, func(c Class) bool { return c.ID == eventData.ClassID })
	if class == nil {
		return ErrApplyEventClassNotFound(event.EventType(), eventData.ClassID)
	}
	a.Version = event.EventVersion()
	class.UpdatedAt = event.EventAt()
	class.NumberOfPupils += eventData.Number
	return nil
}

func (a *SchoolClassAggregate) onNumberOfPupilsDecreased(event domain.Event) error {
	eventData := NumberOfPupilsDecreasedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	class := fp.Find(a.Classes, func(c Class) bool { return c.ID == eventData.ClassID })
	if class == nil {
		return ErrApplyEventClassNotFound(event.EventType(), eventData.ClassID)
	}
	a.Version = event.EventVersion()
	class.UpdatedAt = event.EventAt()
	class.NumberOfPupils -= eventData.Number
	return nil
}
