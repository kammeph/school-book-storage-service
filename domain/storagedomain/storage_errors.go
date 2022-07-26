package storagedomain

import (
	"errors"
	"fmt"
)

var (
	ErrStorageNameNotSet     = errors.New("storage name not set")
	ErrStorageLocationNotSet = errors.New("storage location not set")
)

func ErrStoragesWithIdAlreadyExists(id string) error {
	return fmt.Errorf("storage with ID %s already exists", id)
}

func ErrStorageAlreadyExists(name, location string) error {
	return fmt.Errorf("storage with name %s in location %s already exists", name, location)
}

func ErrStorageIDNotFound(id string) error {
	return fmt.Errorf("storage with ID %s not found", id)
}

func ErrStorageByNameNotFound(name string) error {
	return fmt.Errorf("storage with the name %s does not exist", name)
}

func ErrMultipleStoragesWithNameFound(name string) error {
	return fmt.Errorf("there are more than one storage with the name %s", name)
}
