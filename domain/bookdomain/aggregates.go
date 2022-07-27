package bookdomain

import (
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

type SchoolBookAggregate struct {
	*domain.AggregateModel
	Books []Book
}

func NewSchoolBookAggregate() *SchoolBookAggregate {
	aggregate := &SchoolBookAggregate{}
	model := domain.NewAggregateModel(aggregate.On)
	aggregate.AggregateModel = &model
	return aggregate
}

func (a *SchoolBookAggregate) On(event domain.Event) error {
	switch event.EventType() {
	case BookAdded:
		return a.onBookAdded(event)
	case BookMetaAdjusted:
		return a.onBookMetaAdjusted(event)
	case BookPriceIncreased:
		return a.onBookPriceIncreased(event)
	case BookPriceDecreased:
		return a.onBookPriceDecreased(event)
	default:
		return domain.ErrUnknownEvent(event)
	}
}

func (a *SchoolBookAggregate) onBookAdded(event domain.Event) error {
	var eventData BookAddedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	if fp.Some(a.Books, func(b Book) bool { return b.ID == eventData.BookID }) {
		return ErrApplyEventBookAlreadyExists(event.EventType(), eventData.BookID)
	}
	book := NewBook(
		eventData.BookID,
		eventData.Isbn,
		eventData.Name,
		eventData.Description,
		eventData.Price,
		eventData.Grades,
		event.EventAt(),
	)
	a.Version = event.EventVersion()
	a.Books = append(a.Books, book)
	return nil
}

func (a *SchoolBookAggregate) onBookMetaAdjusted(event domain.Event) error {
	var eventData BookMetaAdjustedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.ID == eventData.BookID })
	if book == nil {
		return ErrApplyEventBookWithIDNotFound(event.EventType(), eventData.BookID)
	}
	a.Version = event.EventVersion()
	book.UpdatedAt = event.EventAt()
	book.Name = eventData.Name
	book.Description = eventData.Description
	book.Grades = eventData.Grades
	return nil
}

func (a *SchoolBookAggregate) onBookPriceIncreased(event domain.Event) error {
	var eventData BookPriceIncreasedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.ID == eventData.BookID })
	if book == nil {
		return ErrApplyEventBookWithIDNotFound(event.EventType(), eventData.BookID)
	}
	a.Version = event.EventVersion()
	book.UpdatedAt = event.EventAt()
	book.Price = eventData.Price
	return nil
}

func (a *SchoolBookAggregate) onBookPriceDecreased(event domain.Event) error {
	var eventData BookPriceDecreasedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	book := fp.Find(a.Books, func(b Book) bool { return b.ID == eventData.BookID })
	if book == nil {
		return ErrApplyEventBookWithIDNotFound(event.EventType(), eventData.BookID)
	}
	a.Version = event.EventVersion()
	book.UpdatedAt = event.EventAt()
	book.Price = eventData.Price
	return nil
}
