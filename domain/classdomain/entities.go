package classdomain

import "time"

type Class struct {
	ID             string
	Grade          int
	Letter         string
	NumberOfPupils int
	DateFrom       time.Time
	DateTo         time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewClass(id string, grade int, letter string, pupils int, from, to, timeStamp time.Time) Class {
	return Class{
		ID:             id,
		Grade:          grade,
		Letter:         letter,
		NumberOfPupils: pupils,
		DateFrom:       from,
		DateTo:         to,
		CreatedAt:      timeStamp,
	}
}
