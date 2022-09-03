package schoolapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
)

type SchoolCommandHandlers struct {
	AddSchoolHandler      *AddSchoolCommandHandler
	DeactiveSchoolHandler *DeactivateSchoolCommandHandler
	RenameSchoolHandler   *RenameSchoolCommandHandler
}

func NewSchoolCommandHandlers(store application.Store, publisher application.EventPublisher) *SchoolCommandHandlers {
	return &SchoolCommandHandlers{
		AddSchoolHandler:      NewAddStorageCommandHandler(store, publisher),
		DeactiveSchoolHandler: NewDeactivateStorageCommandHandler(store, publisher),
		RenameSchoolHandler:   NewRenameStorageCommandHandler(store, publisher),
	}
}

type AddSchoolCommand struct {
	application.CommandModel
	Name string `json:"name"`
}

type AddSchoolCommandHandler struct {
	*application.CommandHandlerModel
}

func NewAddStorageCommandHandler(store application.Store, publisher application.EventPublisher) *AddSchoolCommandHandler {
	return &AddSchoolCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h *AddSchoolCommandHandler) Handle(ctx context.Context, command AddSchoolCommand) (string, error) {
	aggregate := schooldomain.NewSchoolAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return "", err
	}
	schoolID, err := aggregate.AddSchool(command.Name)
	if err != nil {
		return "", err
	}
	if err := h.SaveAndPublish(ctx, aggregate); err != nil {
		return "", err
	}
	return schoolID, nil
}

type DeactivateSchoolCommand struct {
	application.CommandModel
	SchoolID string
	Reason   string
}

type DeactivateSchoolCommandHandler struct {
	*application.CommandHandlerModel
}

func NewDeactivateStorageCommandHandler(store application.Store, publisher application.EventPublisher) *DeactivateSchoolCommandHandler {
	return &DeactivateSchoolCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h *DeactivateSchoolCommandHandler) Handle(ctx context.Context, command DeactivateSchoolCommand) error {
	aggregate := schooldomain.NewSchoolAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.DeactivateSchool(command.SchoolID, command.Reason); err != nil {
		return nil
	}
	if err := h.SaveAndPublish(ctx, aggregate); err != nil {
		return err
	}
	return nil
}

type RenameSchoolCommand struct {
	application.CommandModel
	SchoolID string
	Name     string
	Reason   string
}

type RenameSchoolCommandHandler struct {
	*application.CommandHandlerModel
}

func NewRenameStorageCommandHandler(store application.Store, publisher application.EventPublisher) *RenameSchoolCommandHandler {
	return &RenameSchoolCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h *RenameSchoolCommandHandler) Handle(ctx context.Context, command RenameSchoolCommand) error {
	aggregate := schooldomain.NewSchoolAggregateWithID(command.AggregateID())
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RenameSchool(command.SchoolID, command.Name, command.Reason); err != nil {
		return err
	}
	if err := h.SaveAndPublish(ctx, aggregate); err != nil {
		return err
	}
	return nil
}
