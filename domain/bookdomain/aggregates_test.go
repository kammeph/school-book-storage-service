package bookdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/bookdomain"
	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/stretchr/testify/assert"
)

type UnknownEvent struct {
	domain.EventModel
}

func TestSchoolBookAggregateOn(t *testing.T) {
	tests := []struct {
		name            string
		eventType       string
		eventVersion    int
		eventAt         time.Time
		schoolID        string
		bookID          string
		isbn            string
		bookName        string
		bookDescription string
		price           float64
		grades          []int
		reason          string
		err             error
		expectError     bool
		addDefaultBooks bool
	}{
		{
			name:            "Book created",
			eventType:       bookdomain.BookAdded,
			eventVersion:    4,
			eventAt:         time.Now(),
			schoolID:        "school",
			isbn:            "12345",
			bookID:          "book",
			bookName:        "test book",
			bookDescription: "test description",
			price:           30.00,
			grades:          []int{3, 4},
			err:             nil,
			expectError:     false,
			addDefaultBooks: false,
		},
		{
			name:            "Book created error",
			eventType:       bookdomain.BookAdded,
			eventVersion:    9,
			eventAt:         time.Now(),
			schoolID:        "school",
			bookID:          "book",
			bookName:        "A",
			err:             bookdomain.ErrApplyEventBookAlreadyExists(bookdomain.BookAdded, "book"),
			expectError:     true,
			addDefaultBooks: true,
		},
		{
			name:            "Book meta adjusted",
			eventType:       bookdomain.BookMetaAdjusted,
			eventVersion:    7,
			eventAt:         time.Now(),
			bookID:          "book",
			bookName:        "test book",
			bookDescription: "test description",
			grades:          []int{2, 4},
			err:             nil,
			expectError:     false,
			addDefaultBooks: true,
		},
		{
			name:            "Book meta adjusted Book not found",
			eventType:       bookdomain.BookMetaAdjusted,
			eventVersion:    9,
			eventAt:         time.Now(),
			bookID:          "unkwon",
			err:             bookdomain.ErrApplyEventBookWithIDNotFound(bookdomain.BookMetaAdjusted, "unkwon"),
			expectError:     true,
			addDefaultBooks: false,
		},
		{
			name:            "Book price increased",
			eventType:       bookdomain.BookPriceIncreased,
			eventVersion:    21,
			eventAt:         time.Now(),
			bookID:          "book",
			price:           46,
			err:             nil,
			expectError:     false,
			addDefaultBooks: true,
		},
		{
			name:            "Book price increased error",
			eventType:       bookdomain.BookPriceIncreased,
			eventVersion:    52,
			eventAt:         time.Now(),
			bookID:          "unkwon",
			err:             bookdomain.ErrApplyEventBookWithIDNotFound(bookdomain.BookPriceIncreased, "unkwon"),
			expectError:     true,
			addDefaultBooks: false,
		},
		{
			name:            "Book price decreased",
			eventType:       bookdomain.BookPriceDecreased,
			eventVersion:    21,
			eventAt:         time.Now(),
			bookID:          "book",
			price:           26.32,
			err:             nil,
			expectError:     false,
			addDefaultBooks: true,
		},
		{
			name:            "Book price decreased error",
			eventType:       bookdomain.BookPriceDecreased,
			eventVersion:    52,
			eventAt:         time.Now(),
			bookID:          "unkwon",
			err:             bookdomain.ErrApplyEventBookWithIDNotFound(bookdomain.BookPriceDecreased, "unkwon"),
			expectError:     true,
			addDefaultBooks: false,
		},
		{
			name:            "unknown",
			eventType:       "unknown",
			eventVersion:    52,
			eventAt:         time.Now(),
			err:             domain.ErrUnknownEvent(&UnknownEvent{}),
			expectError:     true,
			addDefaultBooks: false,
		},
	}
	for _, test := range tests {
		var event domain.Event
		t.Run(test.name, func(t *testing.T) {
			switch test.eventType {
			case bookdomain.BookAdded:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := bookdomain.BookAddedEvent{
					test.schoolID,
					test.bookID,
					test.isbn,
					test.bookName,
					test.bookDescription,
					test.price,
					test.grades,
				}
				event.SetJsonData(eventData)
			case bookdomain.BookMetaAdjusted:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := bookdomain.BookMetaAdjustedEvent{
					test.bookID,
					test.bookName,
					test.bookDescription,
					test.grades,
				}
				event.SetJsonData(eventData)
			case bookdomain.BookPriceIncreased:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := bookdomain.BookPriceIncreasedEvent{
					test.bookID,
					test.price,
					test.reason,
				}
				event.SetJsonData(eventData)
			case bookdomain.BookPriceDecreased:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := bookdomain.BookPriceDecreasedEvent{
					test.bookID,
					test.price,
					test.reason,
				}
				event.SetJsonData(eventData)
			default:
				event = &UnknownEvent{}
			}
			aggregate := bookdomain.NewSchoolBookAggregate()
			if test.addDefaultBooks {
				aggregate.Books = append(aggregate.Books, bookdomain.Book{
					ID:          test.bookID,
					Isbn:        test.isbn,
					Name:        "biologie book",
					Description: "book for biologie lessons",
					Price:       59.99,
					Grades:      []int{5},
				})
			}
			err := aggregate.On(event)
			if test.expectError {
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.eventVersion, aggregate.Version)
			Book := fp.Find(aggregate.Books, func(c bookdomain.Book) bool { return c.ID == test.bookID })
			assert.NotNil(t, Book)
			if test.eventType == bookdomain.BookAdded {
				assert.Equal(t, test.eventVersion, aggregate.Version)
				assert.Equal(t, test.bookID, Book.ID)
				assert.Equal(t, test.isbn, Book.Isbn)
				assert.Equal(t, test.bookName, Book.Name)
				assert.Equal(t, test.bookDescription, Book.Description)
				assert.Equal(t, test.price, Book.Price)
				assert.Equal(t, test.grades, Book.Grades)
				assert.Equal(t, test.eventAt, Book.CreatedAt)
			}
			if test.eventType == bookdomain.BookMetaAdjusted {
				assert.Equal(t, test.eventVersion, aggregate.Version)
				assert.Equal(t, test.bookName, Book.Name)
				assert.Equal(t, test.bookDescription, Book.Description)
				assert.Equal(t, test.grades, Book.Grades)
				assert.Equal(t, test.eventAt, Book.UpdatedAt)
			}
			if test.eventType == bookdomain.BookPriceIncreased || test.eventType == bookdomain.BookPriceDecreased {
				assert.Equal(t, test.eventVersion, aggregate.Version)
				assert.Equal(t, test.price, Book.Price)
				assert.Equal(t, test.eventAt, Book.UpdatedAt)
			}
		})
	}
}
