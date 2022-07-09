package storage

import (
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	ID        string
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStorage(id string, timeStamp time.Time) Storage {
	return Storage{ID: id, CreatedAt: timeStamp, UpdatedAt: timeStamp}
}

type School struct {
	ID        string
	UpdatedAt time.Time
	Storages  []Storage
}

func NewSchool() School {
	return School{ID: uuid.New().String(), Storages: []Storage{}}
}

func NewSchoolWithID(id string) School {
	return School{ID: id, Storages: []Storage{}}
}
