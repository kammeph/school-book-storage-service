package classdomain

import (
	"time"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

func (a *SchoolClassAggregate) CreateClass(grade int, letter string, numberOfPupils int, dateFrom, dateTo time.Time) (string, error) {
	if grade < 1 {
		return "", ErrGradeGreaterZero
	}
	if letter == "" {
		return "", ErrLetterNotSet
	}
	if len(letter) > 1 {
		return "", ErrLetterToLong
	}
	if numberOfPupils < 1 {
		return "", ErrNumberOfPupilsGreaterZero
	}
	if dateTo.Sub(dateFrom) < 0 {
		return "", ErrInvalidDates
	}
	if fp.Some(a.Classes, func(c Class) bool {
		return c.Grade == grade &&
			c.Letter == letter &&
			c.DateFrom.Year() == dateFrom.Year() &&
			c.DateTo.Year() == dateTo.Year()
	}) {
		return "", ErrClassAlreadyExists(grade, letter, dateFrom.Year(), dateTo.Year())
	}
	classID := uuid.NewString()
	event, err := NewClassCreated(a, classID, grade, letter, numberOfPupils, dateFrom, dateTo)
	if err != nil {
		return "", err
	}
	if err := a.Apply(event); err != nil {
		return "", err
	}
	return classID, nil
}

func (a *SchoolClassAggregate) IncreaseNumberOfPupils(classID string, number int, reason string) error {
	if number < 1 {
		return ErrIncreasePupilsGreaterZero
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	class := fp.Find(a.Classes, func(c Class) bool { return c.ID == classID })
	if class == nil {
		return ErrClassWithIDNotFound(classID)
	}
	event, err := NewNumberOfPupilsIncreased(a, classID, number, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *SchoolClassAggregate) DecreaseNumberOfPupils(classID string, number int, reason string) error {
	if number < 1 {
		return ErrDecreasePupilsGreaterZero
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	class := fp.Find(a.Classes, func(c Class) bool { return c.ID == classID })
	if class == nil {
		return ErrClassWithIDNotFound(classID)
	}
	event, err := NewNumberOfPupilsDecreased(a, classID, number, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}
