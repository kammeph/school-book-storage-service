package userdomain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/fp"
	"golang.org/x/crypto/bcrypt"
)

func (a *UsersAggregate) RegisterUser(name, password string, locale Locale) error {
	if name == "" {
		return fmt.Errorf("user name not set")
	}
	if password == "" {
		return fmt.Errorf("user password not set")
	}
	if fp.Some(a.Users, func(u UserModel) bool { return u.Name == name && u.Active }) {
		return fmt.Errorf("the user with the name %s already exists", name)
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userID := uuid.NewString()
	event, err := NewUserRegisteredEvent(a, userID, "", name, passwordHash, []Role{User}, locale)
	if err != nil {
		return err
	}
	return a.Apply(event)
}

func (a *UsersAggregate) DeactivateUser(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID not set")
	}
	user := fp.Find(a.Users, func(u UserModel) bool { return u.ID == userID })
	if user == nil {
		return fmt.Errorf("user with ID %s not found", userID)
	}
	if !user.Active {
		return fmt.Errorf("the user %s is already deactivated", user.Name)
	}
	event, err := NewUserDeactivatedEvent(a, userID)
	if err != nil {
		return err
	}
	return a.Apply(event)
}

func (a *UsersAggregate) LoginUser(name, password string) (*UserModel, error) {
	user := fp.Find(a.Users, func(u UserModel) bool { return u.Name == name && u.Active })
	if user == nil {
		return nil, fmt.Errorf("user with name %s not found", name)
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return nil, errors.New("password is incorrect")
	}
	return user, nil
}
