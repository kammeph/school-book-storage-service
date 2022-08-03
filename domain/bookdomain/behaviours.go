package bookdomain

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

func (a *SchoolBookAggregate) AddBook(isbn, name, description string, price float64, grades []int) (string, error) {
	if isbn == "" {
		return "", ErrIsbnNotSet
	}
	if name == "" {
		return "", ErrBookNameNotSet
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.Name == name || b.Isbn == isbn })
	if book != nil {
		return "", ErrBookAlreadyExists(isbn, name)
	}
	bookID := uuid.NewString()
	event, err := NewBookAddedEvent(a, a.AggregateID(), bookID, isbn, name, description, price, grades)
	if err != nil {
		return "", err
	}
	if err := a.Apply(event); err != nil {
		return "", err
	}
	return bookID, nil
}

func (a *SchoolBookAggregate) AdjustBookMeta(bookID, name, description string, grades []int) error {
	if name == "" {
		return ErrBookNameNotSet
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.ID == bookID })
	if book == nil {
		return ErrBookWithIDNotFound(bookID)
	}
	book = fp.Find(a.Books, func(b Book) bool { return b.Name == name })
	if book != nil {
		return ErrBookWithNameAlreadyExists(name)
	}
	event, err := NewBookMetaAdjustedEvent(a, bookID, name, description, grades)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *SchoolBookAggregate) IncreaseBookPrice(bookID string, price float64, reason string) error {
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	if price < 0.0 {
		return ErrBookPriceLessThanZero
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.ID == bookID })
	if book == nil {
		return ErrBookWithIDNotFound(bookID)
	}
	event, err := NewBookPriceIncreasedEvent(a, bookID, reason, price)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *SchoolBookAggregate) DecreaseBookPrice(bookID string, price float64, reason string) error {
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	if price < 0.0 {
		return ErrBookPriceLessThanZero
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.ID == bookID })
	if book == nil {
		return ErrBookWithIDNotFound(bookID)
	}
	event, err := NewBookPriceDecreasedEvent(a, bookID, reason, price)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}
