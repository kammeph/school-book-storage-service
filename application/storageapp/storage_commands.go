package storageapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain/storagedomain"
)

type StorageCommandHandlers struct {
	AddStorageHandler      AddStorageCommandHandler
	RemoveStorageHandler   RemoveStorageCommandHandler
	RenameStorageHandler   RenameStorageCommandHandler
	RelocateStorageHandler RelocateStorageCommandHandler
}

func NewStorageCommandHandlers(store application.Store, publisher application.EventPublisher) StorageCommandHandlers {
	return StorageCommandHandlers{
		AddStorageHandler:      NewAddStorageCommandHandler(store, publisher),
		RemoveStorageHandler:   NewRemoveStorageCommandHandler(store, publisher),
		RenameStorageHandler:   NewRenameStorageCommandHandler(store, publisher),
		RelocateStorageHandler: NewRelocateStorageCommandHandler(store, publisher),
	}
}

type AddStorageCommand struct {
	application.CommandModel
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AddStorageCommandHandler struct {
	*application.CommandHandlerModel
}

func NewAddStorageCommandHandler(store application.Store, publisher application.EventPublisher) AddStorageCommandHandler {
	return AddStorageCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h AddStorageCommandHandler) Handle(ctx context.Context, command AddStorageCommand) (string, error) {
	aggregate := storagedomain.NewSchoolStorageAggregateWithID(command.AggregateID())
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
	application.CommandModel
	StorageID string `json:"storageId"`
	Reason    string `json:"reason"`
}

type RemoveStorageCommandHandler struct {
	*application.CommandHandlerModel
}

func NewRemoveStorageCommandHandler(store application.Store, publisher application.EventPublisher) RemoveStorageCommandHandler {
	return RemoveStorageCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h RemoveStorageCommandHandler) Handle(ctx context.Context, command RemoveStorageCommand) error {
	aggregate := storagedomain.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RemoveStorage(command.StorageID, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type RenameStorageCommand struct {
	application.CommandModel
	StorageID string `json:"storageId"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
}

type RenameStorageCommandHandler struct {
	*application.CommandHandlerModel
}

func NewRenameStorageCommandHandler(store application.Store, publisher application.EventPublisher) RenameStorageCommandHandler {
	return RenameStorageCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h RenameStorageCommandHandler) Handle(ctx context.Context, command RenameStorageCommand) error {
	aggregate := storagedomain.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RenameStorage(command.StorageID, command.Name, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type RelocateStorageCommand struct {
	application.CommandModel
	StorageID string `json:"storageId"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
}

type RelocateStorageCommandHandler struct {
	*application.CommandHandlerModel
}

func NewRelocateStorageCommandHandler(store application.Store, publisher application.EventPublisher) RelocateStorageCommandHandler {
	return RelocateStorageCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h RelocateStorageCommandHandler) Handle(ctx context.Context, command RelocateStorageCommand) error {
	aggregate := storagedomain.NewSchoolStorageAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RelocateStorage(command.StorageID, command.Location, command.Reason); err != nil {
		return err
	}
	return h.SaveAndPublish(ctx, aggregate)
}
