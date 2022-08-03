package bookdomain_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain/bookdomain"
	"github.com/stretchr/testify/assert"
)

func TestNewBookAddedEvent(t *testing.T) {
	schoolID := "school"
	bookID := "masterBook"
	isbn := "12345"
	name := "English Book"
	description := "book for english lessons"
	price := 20.99
	grades := []int{2, 3}
	aggregate := bookdomain.NewSchoolBookAggregate()
	aggregate.Version = 4
	event, err := bookdomain.NewBookAddedEvent(aggregate, schoolID, bookID, isbn, name, description, price, grades)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, bookdomain.BookAdded, event.EventType())
	assert.Equal(t, 5, event.EventVersion())
	assert.NotZero(t, event.EventAt())
	assert.NotEmpty(t, event)
}

func TestNewBookMetaAdjustedEvent(t *testing.T) {
	bookID := "masterbook"
	name := "Math Book"
	description := "book for math lessons"
	grades := []int{1, 2}
	aggregate := bookdomain.NewSchoolBookAggregate()
	aggregate.Version = 10
	event, err := bookdomain.NewBookMetaAdjustedEvent(aggregate, bookID, name, description, grades)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, bookdomain.BookMetaAdjusted, event.EventType())
	assert.Equal(t, 11, event.EventVersion())
	assert.NotZero(t, event.EventAt())
	assert.NotEmpty(t, event.EventData())
}

func TestNewBookPriceIncreasedEvent(t *testing.T) {
	bookID := "masterbook"
	price := 19.99
	reason := "test"
	aggregate := bookdomain.NewSchoolBookAggregate()
	aggregate.Version = 7
	event, err := bookdomain.NewBookPriceIncreasedEvent(aggregate, bookID, reason, price)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, bookdomain.BookPriceIncreased, event.EventType())
	assert.Equal(t, 8, event.EventVersion())
	assert.NotZero(t, event.EventAt())
	assert.NotEmpty(t, event.EventData())
}

func TestNewBookPriceDecreasedEvent(t *testing.T) {
	bookID := "masterbook"
	price := 19.99
	reason := "test"
	aggregate := bookdomain.NewSchoolBookAggregate()
	aggregate.Version = 7
	event, err := bookdomain.NewBookPriceDecreasedEvent(aggregate, bookID, reason, price)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, bookdomain.BookPriceDecreased, event.EventType())
	assert.Equal(t, 8, event.EventVersion())
	assert.NotZero(t, event.EventAt())
	assert.NotEmpty(t, event.EventData())
}
