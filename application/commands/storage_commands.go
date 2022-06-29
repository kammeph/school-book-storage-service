package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/domain/aggregates"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
)

type AddStorage struct {
	common.CommandModel
	Name     string
	Location string
}

type RemoveStorage struct {
	common.CommandModel
	Reason string
}

type SetStorageName struct {
	common.CommandModel
	Name   string
	Reason string
}

type SetStorageLocation struct {
	common.CommandModel
	Location string
	Reason   string
}

type StorageCommandHandler struct {
}

func (h *StorageCommandHandler) Apply(ctx context.Context, aggregate domain.Aggregate, command common.Command) ([]domain.Event, error) {
	storage, ok := aggregate.(*aggregates.StorageAggregate)
	if !ok {
		return nil, fmt.Errorf("Incorrect type for aggregate: %t", aggregate)
	}
	switch c := command.(type) {
	case AddStorage:
		return storage.AddStorage(c.ID, c.Name, c.Location)
	case RemoveStorage:
		return storage.RemoveStorage(c.Reason)
	case SetStorageName:
		return storage.SetStorageName(c.Name, c.Reason)
	case SetStorageLocation:
		return storage.SetStorageLocation(c.Location, c.Reason)
	default:
		return nil, errors.New("Unapplied command")
	}
}
