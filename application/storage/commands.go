package storage

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageCommandHandlers struct {
	AddStorageHandler         AddStorageCommandHandler
	RemoveStorageHandler      RemoveStorageCommandHandler
	SetStorageNameHandler     SetStorageNameCommandHandler
	SetStorageLocationHandler SetStorageLocationCommandHandler
}

func NewStorageCommandHandlers(repository *common.Repository) StorageCommandHandlers {
	return StorageCommandHandlers{
		AddStorageHandler:         NewAddStorageCommandHandler(repository),
		RemoveStorageHandler:      NewRemoveStorageCommandHandler(repository),
		SetStorageNameHandler:     NewSetStorageNameCommandHandler(repository),
		SetStorageLocationHandler: NewSetStorageLocationCommandHandler(repository),
	}
}

type AddStorage struct {
	common.CommandModel
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AddStorageCommandHandler struct {
	repository *common.Repository
}

func NewAddStorageCommandHandler(repository *common.Repository) AddStorageCommandHandler {
	return AddStorageCommandHandler{repository: repository}
}

func (h AddStorageCommandHandler) Handle(ctx context.Context, command AddStorage) (StorageIDDto, error) {
	aggregate, err := h.repository.Load(ctx, command.AggregateID())
	if err != nil {
		return StorageIDDto{}, err
	}
	school, ok := aggregate.(*storage.StorageAggregateRoot)
	if !ok {
		return StorageIDDto{}, IncorrectAggregateTypeError(aggregate)
	}
	storageID, err := school.AddStorage(command.Name, command.Location)
	if err != nil {
		return StorageIDDto{}, err
	}
	err = h.repository.Save(ctx, aggregate)
	return StorageIDDto{StorageID: storageID}, err
}

type RemoveStorage struct {
	common.CommandModel
	StorageID string `json:"storageId"`
	Reason    string `json:"reason"`
}

type RemoveStorageCommandHandler struct {
	repository *common.Repository
}

func NewRemoveStorageCommandHandler(repository *common.Repository) RemoveStorageCommandHandler {
	return RemoveStorageCommandHandler{repository: repository}
}

func (h RemoveStorageCommandHandler) Handle(ctx context.Context, command RemoveStorage) error {
	aggregate, err := h.repository.Load(ctx, command.AggregateID())
	if err != nil {
		return err
	}
	school, ok := aggregate.(*storage.StorageAggregateRoot)
	if !ok {
		return IncorrectAggregateTypeError(aggregate)
	}
	err = school.RemoveStorage(command.StorageID, command.Reason)
	if err != nil {
		return err
	}
	return h.repository.Save(ctx, aggregate)
}

type SetStorageName struct {
	common.CommandModel
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
}

type SetStorageNameCommandHandler struct {
	repository *common.Repository
}

func NewSetStorageNameCommandHandler(repository *common.Repository) SetStorageNameCommandHandler {
	return SetStorageNameCommandHandler{repository: repository}
}

func (h SetStorageNameCommandHandler) Handle(ctx context.Context, command SetStorageName) error {
	aggregate, err := h.repository.Load(ctx, command.AggregateID())
	if err != nil {
		return err
	}
	school, ok := aggregate.(*storage.StorageAggregateRoot)
	if !ok {
		return IncorrectAggregateTypeError(aggregate)
	}
	err = school.RenameStorage(command.StorageID, command.Name, command.Reason)
	if err != nil {
		return err
	}
	return h.repository.Save(ctx, aggregate)
}

type SetStorageLocation struct {
	common.CommandModel
	StorageID string `json:"storageId"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
}

type SetStorageLocationCommandHandler struct {
	repository *common.Repository
}

func NewSetStorageLocationCommandHandler(repository *common.Repository) SetStorageLocationCommandHandler {
	return SetStorageLocationCommandHandler{repository: repository}
}

func (h SetStorageLocationCommandHandler) Handle(ctx context.Context, command SetStorageLocation) error {
	aggregate, err := h.repository.Load(ctx, command.AggregateID())
	if err != nil {
		return err
	}
	school, ok := aggregate.(*storage.StorageAggregateRoot)
	if !ok {
		return IncorrectAggregateTypeError(aggregate)
	}
	err = school.RelocateStorage(command.StorageID, command.Location, command.Reason)
	if err != nil {
		return err
	}
	return h.repository.Save(ctx, aggregate)
}
