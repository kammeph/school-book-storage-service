package domain

import (
	"errors"
	"fmt"
)

var ErrSchoolNameNotSet = errors.New("school name not set")

func ErrSchoolWithIdAlreadyExists(id string) error {
	return fmt.Errorf("school with ID %s already exists", id)
}

func ErrSchoolAlreadyExists(name string) error {
	return fmt.Errorf("school with name %s already exists", name)
}

func ErrSchoolWithIDNotFound(id string) error {
	return fmt.Errorf("school with ID %s not found", id)
}

func ErrSchoolWithNameNotFound(name string) error {
	return fmt.Errorf("school with the name %s does not exist", name)
}

func ErrMultipleSchoolsWithNameFound(name string) error {
	return fmt.Errorf("there are more than one schools with the name %s", name)
}
