package bookdomain

import (
	"errors"
	"fmt"
)

var (
	ErrIsbnNotSet            = errors.New("ISBN not set")
	ErrBookNameNotSet        = errors.New("Book name not set")
	ErrBookDescriptionNotSet = errors.New("Book description not set")
	ErrBookPriceLessThanZero = errors.New("Book price is less than zero")
)

func ErrApplyEventBookAlreadyExists(eventType, id string) error {
	return fmt.Errorf("can not apply %s: Book with ID %s already exists", eventType, id)
}

func ErrApplyEventBookWithIDNotFound(eventType, id string) error {
	return fmt.Errorf("can not apply %s: Book with ID %s not found", eventType, id)
}

func ErrBookWithIDNotFound(bookID string) error {
	return fmt.Errorf("book with ID %s not found", bookID)
}

func ErrBookAlreadyExists(isbn, name string) error {
	return fmt.Errorf("there is already a book with the ISBN %s or name %s", isbn, name)
}

func ErrBookWithNameAlreadyExists(name string) error {
	return fmt.Errorf("there is already a book with the name %s", name)
}
