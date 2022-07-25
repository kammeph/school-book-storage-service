package application

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
)

type StorageCommandHandlers struct {
	AddStorageHandler      AddStorageCommandHandler
	RemoveStorageHandler   RemoveStorageCommandHandler
	RenameStorageHandler   RenameStorageCommandHandler
	RelocateStorageHandler RelocateStorageCommandHandler
}

func NewStorageCommandHandlers(store Store, publisher EventPublisher) StorageCommandHandlers {
	return StorageCommandHandlers{
		AddStorageHandler:      NewAddStorageCommandHandler(store, publisher),
		RemoveStorageHandler:   NewRemoveStorageCommandHandler(store, publisher),
		RenameStorageHandler:   NewRenameStorageCommandHandler(store, publisher),
		RelocateStorageHandler: NewRelocateStorageCommandHandler(store, publisher),
	}
}

type AddStorageCommand struct {
	CommandModel
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AddStorageCommandHandler struct {
	*CommandHandlerModel
}

func NewAddStorageCommandHandler(store Store, publisher EventPublisher) AddStorageCommandHandler {
	return AddStorageCommandHandler{NewCommandHandlerModel(store, publisher)}
}

func (h AddStorageCommandHandler) Handle(ctx context.Context, command AddStorageCommand) (string, error) {
	aggregate := domain.NewSchoolStorageAggregateWithID(command.AggregateID())
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
	CommandModel
	StorageID string `json:"storageId"`
	Reason    string `json:"reason"`
}

type RemoveStorageCommandHandler struct {
	*CommandHandlerModel
}

func NewRemoveStorageCommandHandler(store Store, publisher EventPublisher) RemoveStorageCommandHandler {
	return RemoveStorageCommandHandler{NewCommandHandlerModel(store, publisher)}
}

func (h RemoveStorageCommandHandler) Handle(ctx context.Context, command RemoveStorageCommand) error {
	aggregate := domain.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RemoveStorage(command.StorageID, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type RenameStorageCommand struct {
	CommandModel
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
}

type RenameStorageCommandHandler struct {
	*CommandHandlerModel
}

func NewRenameStorageCommandHandler(store Store, publisher EventPublisher) RenameStorageCommandHandler {
	return RenameStorageCommandHandler{NewCommandHandlerModel(store, publisher)}
}

func (h RenameStorageCommandHandler) Handle(ctx context.Context, command RenameStorageCommand) error {
	aggregate := domain.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RenameStorage(command.StorageID, command.Name, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type RelocateStorageCommand struct {
	CommandModel
	StorageID string `json:"storageId"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
}

type RelocateStorageCommandHandler struct {
	*CommandHandlerModel
}

func NewRelocateStorageCommandHandler(store Store, publisher EventPublisher) RelocateStorageCommandHandler {
	return RelocateStorageCommandHandler{NewCommandHandlerModel(store, publisher)}
}

func (h RelocateStorageCommandHandler) Handle(ctx context.Context, command RelocateStorageCommand) error {
	aggregate := domain.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RelocateStorage(command.StorageID, command.Location, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}
