package classdomain

import "time"

type Class struct {
	ID        string
	Grade     int
	Letter    string
	YearFrom  int
	YearTo    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewClass(id string, grade int, letter string, from, to int, timeStamp time.Time) Class {
	return Class{
		ID:        id,
		Grade:     grade,
		Letter:    letter,
		YearFrom:  from,
		YearTo:    to,
		CreatedAt: timeStamp,
	}
}
