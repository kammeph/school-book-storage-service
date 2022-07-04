package storage

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func StorageIDNotFoundError(id uuid.UUID) error {
	return fmt.Errorf("Storage with ID %s not found", id)
}

func UnknownEventError(event common.Event) error {
	return fmt.Errorf("Unhandled event %v", event)
}
