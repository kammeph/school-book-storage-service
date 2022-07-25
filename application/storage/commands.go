package storage

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/domain/storage"
)

type StorageCommandHandlers struct {
	AddStorageHandler      AddStorageCommandHandler
	RemoveStorageHandler   RemoveStorageCommandHandler
	RenameStorageHandler   RenameStorageCommandHandler
	RelocateStorageHandler RelocateStorageCommandHandler
}

func NewStorageCommandHandlers(store common.Store, publisher common.EventPublisher) StorageCommandHandlers {
	return StorageCommandHandlers{
		AddStorageHandler:      NewAddStorageCommandHandler(store, publisher),
		RemoveStorageHandler:   NewRemoveStorageCommandHandler(store, publisher),
		RenameStorageHandler:   NewRenameStorageCommandHandler(store, publisher),
		RelocateStorageHandler: NewRelocateStorageCommandHandler(store, publisher),
	}
}

type AddStorageCommand struct {
	common.CommandModel
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AddStorageCommandHandler struct {
	*common.CommandHandlerModel
}

func NewAddStorageCommandHandler(store common.Store, publisher common.EventPublisher) AddStorageCommandHandler {
	return AddStorageCommandHandler{common.NewCommandHandlerModel(store, publisher)}
}

func (h AddStorageCommandHandler) Handle(ctx context.Context, command AddStorageCommand) (string, error) {
	aggregate := storage.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return "", err
	}
	storageID, err := aggregate.AddStorage(command.Name, command.Location)
	if err != nil {
		return "", err
	}
	if err := h.SaveAndPublish(ctx, aggregate); err != nil {
		return "", err
	}
	return storageID, err
}

type RemoveStorageCommand struct {
	common.CommandModel
	StorageID string `json:"storageId"`
	Reason    string `json:"reason"`
}

type RemoveStorageCommandHandler struct {
	*common.CommandHandlerModel
}

func NewRemoveStorageCommandHandler(store common.Store, publisher common.EventPublisher) RemoveStorageCommandHandler {
	return RemoveStorageCommandHandler{common.NewCommandHandlerModel(store, publisher)}
}

func (h RemoveStorageCommandHandler) Handle(ctx context.Context, command RemoveStorageCommand) error {
	aggregate := storage.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RemoveStorage(command.StorageID, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type RenameStorageCommand struct {
	common.CommandModel
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
}

type RenameStorageCommandHandler struct {
	*common.CommandHandlerModel
}

func NewRenameStorageCommandHandler(store common.Store, publisher common.EventPublisher) RenameStorageCommandHandler {
	return RenameStorageCommandHandler{common.NewCommandHandlerModel(store, publisher)}
}

func (h RenameStorageCommandHandler) Handle(ctx context.Context, command RenameStorageCommand) error {
	aggregate := storage.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RenameStorage(command.StorageID, command.Name, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type RelocateStorageCommand struct {
	common.CommandModel
	StorageID string `json:"storageId"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
}

type RelocateStorageCommandHandler struct {
	*common.CommandHandlerModel
}

func NewRelocateStorageCommandHandler(store common.Store, publisher common.EventPublisher) RelocateStorageCommandHandler {
	return RelocateStorageCommandHandler{common.NewCommandHandlerModel(store, publisher)}
}

func (h RelocateStorageCommandHandler) Handle(ctx context.Context, command RelocateStorageCommand) error {
	aggregate := storage.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RelocateStorage(command.StorageID, command.Location, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}
