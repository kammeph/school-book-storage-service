package entities

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
	Removed   bool
}

func NewStorage(id uuid.UUID, timeStamp time.Time) Storage {
	return Storage{ID: id, CreatedAt: timeStamp, UpdatedAt: timeStamp}
}
