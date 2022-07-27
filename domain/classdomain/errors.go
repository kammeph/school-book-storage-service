package classdomain

import (
	"errors"
	"fmt"
)

var (
	ErrGradeGreaterZero          = errors.New("the grade must be greater than zero")
	ErrLetterNotSet              = errors.New("letter is not set")
	ErrLetterToLong              = errors.New("the class letter should only have on place")
	ErrNumberOfPupilsGreaterZero = errors.New("the number of pupils must be greater than zero")
	ErrInvalidDates              = errors.New("the end date must be greater the the start date")
	ErrIncreasePupilsGreaterZero = errors.New("the number of pupils should at minimum increased by one")
	ErrDecreasePupilsGreaterZero = errors.New("the number of pupils should at minimum decreased by one")
)

func ErrApplyEventClassAlreadyExists(eventType, classID string) error {
	return fmt.Errorf("can not apply %s: Class with ID %s already exists", eventType, classID)
}

func ErrApplyEventClassNotFound(eventType, classID string) error {
	return fmt.Errorf("can not apply %s: Class with ID %s not found", eventType, classID)
}

func ErrClassAlreadyExists(grade int, letter string, yearFrom, yearTo int) error {
	return fmt.Errorf("there is already a similar class %d%s (%d-%d)", grade, letter, yearFrom, yearTo)
}

func ErrClassWithIDNotFound(id string) error {
	return fmt.Errorf("class with ID %s not found", id)
}
