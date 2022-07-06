package storage

import (
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	ID        uuid.UUID
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStorage(id uuid.UUID, timeStamp time.Time) Storage {
	return Storage{ID: id, CreatedAt: timeStamp, UpdatedAt: timeStamp}
}

type School struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Storages  []Storage
}

func NewSchool() School {
	return School{ID: uuid.New(), Storages: []Storage{}}
}

func NewSchoolWithID(id uuid.UUID) School {
	return School{ID: id, Storages: []Storage{}}
}
