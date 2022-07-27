package classdomain

import (
	"time"

	"github.com/kammeph/school-book-storage-service/domain"
)

var (
	ClassCreated            = "CLASS_CREATED"
	NumberOfPupilsIncreased = "NUMBER_OF_PUPILS_INCREASED"
	NumberOfPupilsDecreased = "NUMBER_OF_PUPILS_DECREASED"
)

type ClassCreatedEvent struct {
	SchoolID       string    `json:"schoolId"`
	ClassID        string    `json:"classId"`
	Grade          int       `json:"grade"`
	Letter         string    `json:"letter"`
	NumberOfPupils int       `json:"numberOfPupils"`
	DateFrom       time.Time `json:"yearFrom"`
	DateTo         time.Time `json:"yearTo"`
}

func NewClassCreated(
	aggregate *SchoolClassAggregate,
	classID string,
	grade int,
	letter string,
	numberOfPupils int,
	dateFrom, dateTo time.Time) (domain.Event, error) {
	eventData := ClassCreatedEvent{
		SchoolID: aggregate.AggregateID(),
		ClassID:  classID,
		Grade:    grade,
		Letter:   letter,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	event := domain.NewEvent(aggregate, ClassCreated)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type NumberOfPupilsIncreasedEvent struct {
	ClassID string
	Number  int
	Reason  string
}

func NewNumberOfPupilsIncreased(aggregate *SchoolClassAggregate, classID string, number int, reason string) (domain.Event, error) {
	eventData := NumberOfPupilsIncreasedEvent{
		ClassID: classID,
		Number:  number,
		Reason:  reason,
	}
	event := domain.NewEvent(aggregate, NumberOfPupilsIncreased)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type NumberOfPupilsDecreasedEvent struct {
	ClassID string
	Number  int
	Reason  string
}

func NewNumberOfPupilsDecreased(aggregate *SchoolClassAggregate, classID string, number int, reason string) (domain.Event, error) {
	eventData := NumberOfPupilsDecreasedEvent{
		ClassID: classID,
		Number:  number,
		Reason:  reason,
	}
	event := domain.NewEvent(aggregate, NumberOfPupilsDecreased)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}
