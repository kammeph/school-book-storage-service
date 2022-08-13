package userapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
)

type UserCommandHandlers struct {
	RegisterUserHandler RegisterUserCommandHandler
	LoginUserHandler    LoginUserCommandHandler
}

func NewUsersCommandHandlers(store application.Store, publisher application.EventPublisher) UserCommandHandlers {
	return UserCommandHandlers{
		NewRegisterUserCommandHandler(store, publisher),
		NewLoginUserCommandHandler(store, publisher),
	}
}

type RegisterUserCommand struct {
	application.CommandModel
	Name     string
	Password string
}

type RegisterUserCommandHandler struct {
	*application.CommandHandlerModel
}

func NewRegisterUserCommandHandler(store application.Store, publisher application.EventPublisher) RegisterUserCommandHandler {
	return RegisterUserCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h RegisterUserCommandHandler) Handle(ctx context.Context, command RegisterUserCommand) error {
	aggregate := userdomain.NewUsersAggregateWithID(command.ID)
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return err
	}
	if err := aggregate.RegisterUser(command.Name, command.Password); err != nil {
		return nil
	}
	return h.SaveAndPublish(ctx, aggregate)
}

type LoginUserCommand struct {
	application.CommandModel
	Name     string
	Password string
}

type LoginUserCommandHandler struct {
	*application.CommandHandlerModel
}

func NewLoginUserCommandHandler(store application.Store, publisher application.EventPublisher) LoginUserCommandHandler {
	return LoginUserCommandHandler{application.NewCommandHandlerModel(store, publisher)}
}

func (h LoginUserCommandHandler) Handle(ctx context.Context, command LoginUserCommand) (*userdomain.UserModel, error) {
	aggregate := userdomain.NewUsersAggregateWithID(command.ID)
	if err := h.LoadAggregate(ctx, aggregate); err != nil {
		return nil, err
	}
	return aggregate.LoginUser(command.Name, command.Password)
}
