package userapp

import (
	"context"
	"fmt"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain/userdomain"
	"github.com/kammeph/school-book-storage-service/fp"
)

type UserQueryHandlers struct {
	GetUserByIDHandler GetUserByIDQueryHandler
}

func NewUserQueryHandlers(store application.Store) UserQueryHandlers {
	return UserQueryHandlers{NewGetUserByIDQueryHandler(store)}
}

type GetUserByIDQuery struct {
	application.QueryModel
	UserID string
}

func NewGetUserByIDQuery(id, userID string) GetUserByIDQuery {
	return GetUserByIDQuery{
		QueryModel: application.QueryModel{ID: id},
		UserID:     userID,
	}
}

type GetUserByIDQueryHandler struct {
	store application.Store
}

func NewGetUserByIDQueryHandler(store application.Store) GetUserByIDQueryHandler {
	return GetUserByIDQueryHandler{store}
}

func (h *GetUserByIDQueryHandler) Handle(ctx context.Context, query GetUserByIDQuery) (*userdomain.UserModel, error) {
	aggregate := userdomain.NewUsersAggregateWithID(query.ID)
	events, err := h.store.Load(ctx, aggregate.AggregateID())
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if err := aggregate.On(event); err != nil {
			return nil, err
		}
	}
	user := fp.Find(aggregate.Users, func(u userdomain.UserModel) bool { return u.ID == query.UserID && u.Active })
	if user == nil {
		return nil, fmt.Errorf("user with ID %s not found", query.UserID)
	}
	return user, nil
}
