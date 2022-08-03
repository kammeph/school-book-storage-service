package bookdomain

import "time"

type Book struct {
	ID          string
	Isbn        string
	Name        string
	Description string
	Price       float64
	Grades      []int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBook(id, isbn, name, description string, price float64, grades []int, timestamp time.Time) Book {
	return Book{
		ID:          id,
		Isbn:        isbn,
		Name:        name,
		Description: description,
		Price:       price,
		Grades:      grades,
		CreatedAt:   timestamp,
	}
}
