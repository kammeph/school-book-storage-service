package userdomain

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/fp"
	"golang.org/x/crypto/bcrypt"
)

func (a *UsersAggregate) RegisterUser(name, password string) error {
	if name == "" {
		return fmt.Errorf("user name not set")
	}
	if password == "" {
		return fmt.Errorf("user password not set")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userID := uuid.NewString()
	event, err := NewUserRegisteredEvent(a, userID, name, passwordHash, []Role{User})
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
