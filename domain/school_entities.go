package domain

import "time"

type School struct {
	ID        string
	Name      string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSchool(id, name string, timeStamp time.Time) School {
	return School{
		ID:        id,
		Name:      name,
		Active:    true,
		CreatedAt: timeStamp,
	}
}
