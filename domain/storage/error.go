package storage

import (
	"errors"
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

var (
	StorageNameNotSetError     = errors.New("Storage name not set")
	StorageLocationNotSetError = errors.New("Storage location not set")
	ReasonNotSpecifiedError    = errors.New("No reason specified")
)

func StorageAlreadyExistsError(name, location string) error {
	return fmt.Errorf("Storage with name %s in location %s already exists", name, location)
}

func StorageIDNotFoundError(id string) error {
	return fmt.Errorf("Storage with ID %s not found", id)
}

func StorageByNameNotFoundError(name string) error {
	return fmt.Errorf("Storage with the name %s does not exist", name)
}

func MultipleStoragesWithNameFoundError(name string) error {
	return fmt.Errorf("There are more than one storage with the name %s", name)
}

func UnknownEventError(event common.Event) error {
	return fmt.Errorf("Unhandled event %T", event)
}
