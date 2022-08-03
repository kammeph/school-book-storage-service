package bookdomain

import "github.com/kammeph/school-book-storage-service/domain"

var (
	BookAdded          = "BOOK_ADDED"
	BookMetaAdjusted   = "BOOK_META_ADJUSTED"
	BookPriceIncreased = "BOOK_PRICE_INCREASED"
	BookPriceDecreased = "BOOK_PRICE_DECREASED"
)

type BookAddedEvent struct {
	SchoolID    string
	BookID      string
	Isbn        string
	Name        string
	Description string
	Price       float64
	Grades      []int
}

func NewBookAddedEvent(
	aggregate *SchoolBookAggregate,
	schoolID, bookID, isbn, name, description string,
	price float64,
	grades []int) (domain.Event, error) {
	eventData := BookAddedEvent{schoolID, bookID, isbn, name, description, price, grades}
	event := domain.NewEvent(aggregate, BookAdded)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type BookMetaAdjustedEvent struct {
	BookID      string
	Name        string
	Description string
	Grades      []int
}

func NewBookMetaAdjustedEvent(aggregate *SchoolBookAggregate, bookID, name, description string, grades []int) (domain.Event, error) {
	eventData := BookMetaAdjustedEvent{bookID, name, description, grades}
	event := domain.NewEvent(aggregate, BookMetaAdjusted)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type BookPriceIncreasedEvent struct {
	BookID string
	Price  float64
	Reason string
}

func NewBookPriceIncreasedEvent(aggregate *SchoolBookAggregate, bookID, reason string, price float64) (domain.Event, error) {
	eventData := BookPriceIncreasedEvent{bookID, price, reason}
	event := domain.NewEvent(aggregate, BookPriceIncreased)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}

type BookPriceDecreasedEvent struct {
	BookID string
	Price  float64
	Reason string
}

func NewBookPriceDecreasedEvent(aggregate *SchoolBookAggregate, bookID, reason string, price float64) (domain.Event, error) {
	eventData := BookPriceDecreasedEvent{bookID, price, reason}
	event := domain.NewEvent(aggregate, BookPriceDecreased)
	if err := event.SetJsonData(eventData); err != nil {
		return nil, err
	}
	return event, nil
}
