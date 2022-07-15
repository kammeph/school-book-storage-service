package storage

import (
	"time"
)

type Storage struct {
	ID        string
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStorage(id, name, location string, timeStamp time.Time) Storage {
	return Storage{
		ID:        id,
		Name:      name,
		Location:  location,
		CreatedAt: timeStamp,
	}
}
