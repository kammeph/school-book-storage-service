package userdomain

import "github.com/kammeph/school-book-storage-service/domain"

const (
	UserRegistered  = "USER_REGISTERED"
	UserLoggedIn    = "USER_LOGGED_IN"
	UserLoggedOut   = "USER_LOGGED_OUT"
	UserDeactivated = "USER_DEACTIVATED"
)

type UserRegisteredEventData struct {
	UserID       string `json:"userId"`
	SchoolID     string `json:"schoolId"`
	Name         string `json:"name"`
	PasswordHash []byte `json:"passwordHash"`
	Roles        []Role `json:"roles"`
	Locale       Locale `json:"locale"`
}

func NewUserRegisteredEvent(
	aggregate domain.Aggregate,
	userID, schoolID, name string,
	passwordHash []byte,
	roles []Role,
	locale Locale,
) (domain.Event, error) {
	eventData := UserRegisteredEventData{
		UserID:       userID,
		SchoolID:     schoolID,
		Name:         name,
		PasswordHash: passwordHash,
		Roles:        roles,
		Locale:       locale,
	}
	event := domain.NewEvent(aggregate, UserRegistered)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type UserLoggedInEventData struct {
	UserID string
}

func NewUserLoggedInEvent(aggregate domain.Aggregate, userID string) (domain.Event, error) {
	eventData := UserLoggedInEventData{
		UserID: userID,
	}
	event := domain.NewEvent(aggregate, UserRegistered)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type UserLoggedOutEventData struct {
	UserID string
}

func NewUserLoggedOutEvent(aggregate domain.Aggregate, userID string) (domain.Event, error) {
	eventData := UserLoggedOutEventData{
		UserID: userID,
	}
	event := domain.NewEvent(aggregate, UserRegistered)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type UserDeactivatedEventData struct {
	UserID string
}

func NewUserDeactivatedEvent(aggregate domain.Aggregate, userID string) (domain.Event, error) {
	eventData := UserDeactivatedEventData{
		UserID: userID,
	}
	event := domain.NewEvent(aggregate, UserRegistered)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}
