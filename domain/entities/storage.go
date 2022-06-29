package entities

import (
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	ID        uuid.UUID
	Version   int
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Removed   bool
}
