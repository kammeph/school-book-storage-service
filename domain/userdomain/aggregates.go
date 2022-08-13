package userdomain

import (
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

type UsersAggregate struct {
	*domain.AggregateModel
	Users []UserModel
}

func NewUsersAggregate() *UsersAggregate {
	aggregate := &UsersAggregate{
		Users: []UserModel{},
	}
	model := domain.NewAggregateModel(aggregate.On)
	aggregate.AggregateModel = &model
	return aggregate
}

func NewUsersAggregateWithID(id string) *UsersAggregate {
	aggregate := NewUsersAggregate()
	aggregate.ID = id
	return aggregate
}

func (a *UsersAggregate) On(event domain.Event) error {
	switch event.EventType() {
	case UserRegistered:
		return a.onUserRegistered(event)
	case UserDeactivated:
		return a.onUserDeactivated(event)
	default:
		return domain.ErrUnknownEvent(event)
	}
}

func (a *UsersAggregate) onUserRegistered(event domain.Event) error {
	eventData := UserRegisteredEventData{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	user := NewUser(
		eventData.UserID,
		eventData.SchoolID,
		eventData.Name,
		eventData.PasswordHash,
		eventData.Roles,
		eventData.Locale)
	a.Version = event.EventVersion()
	a.Users = append(a.Users, user)
	return nil
}

func (a *UsersAggregate) onUserDeactivated(event domain.Event) error {
	eventData := UserDeactivatedEventData{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	user := fp.Find(a.Users, func(u UserModel) bool { return u.ID == eventData.UserID })
	if user == nil {
		return fmt.Errorf("user with ID %s not found", eventData.UserID)
	}
	a.Version = event.EventVersion()
	user.Active = false
	return nil
}
