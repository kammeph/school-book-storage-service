package bookdomain_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/bookdomain"
	"github.com/stretchr/testify/assert"
)

func initSchoolBookAggregate(books []bookdomain.Book) *bookdomain.SchoolBookAggregate {
	aggregate := bookdomain.NewSchoolBookAggregate()
	aggregate.Books = books
	return aggregate
}

func TestAddBook(t *testing.T) {
	tests := []struct {
		name        string
		books       []bookdomain.Book
		isbn        string
		bookName    string
		description string
		price       float64
		grades      []int
		err         error
		expectError bool
	}{
		{
			name:        "add Book",
			books:       []bookdomain.Book{},
			isbn:        "123456",
			bookName:    "English Book",
			description: "Book for english lessons",
			price:       29.99,
			grades:      []int{1},
			err:         nil,
			expectError: false,
		},
		{
			name:        "add Book isbn not set",
			books:       []bookdomain.Book{},
			isbn:        "",
			bookName:    "English Book",
			description: "Book for english lessons",
			price:       29.99,
			grades:      []int{1},
			err:         bookdomain.ErrIsbnNotSet,
			expectError: true,
		},
		{
			name:        "add Book name not set",
			books:       []bookdomain.Book{},
			isbn:        "123456",
			bookName:    "",
			description: "Book for english lessons",
			price:       29.99,
			grades:      []int{1},
			err:         bookdomain.ErrBookNameNotSet,
			expectError: true,
		},
		{
			name: "book with isbn exists",
			books: []bookdomain.Book{
				{
					Isbn: "123456",
				},
			},
			isbn:        "123456",
			bookName:    "English Book",
			description: "Book for english lessons",
			price:       29.99,
			grades:      []int{1},
			err:         bookdomain.ErrBookAlreadyExists("123456", "English Book"),
			expectError: true,
		},
		{
			name: "book with name exists",
			books: []bookdomain.Book{
				{
					Name: "English Book",
				},
			},
			isbn:        "123456",
			bookName:    "English Book",
			description: "Book for english lessons",
			price:       29.99,
			grades:      []int{1},
			err:         bookdomain.ErrBookAlreadyExists("123456", "English Book"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolBookAggregate(test.books)
			BookID, err := aggregate.AddBook(test.isbn, test.bookName, test.description, test.price, test.grades)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			assert.NotEqual(t, "", BookID)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, bookdomain.BookAdded, event.EventType())
		})
	}
}

func TestAdjustBookMeta(t *testing.T) {
	tests := []struct {
		name        string
		books       []bookdomain.Book
		BookID      string
		bookName    string
		description string
		grades      []int
		err         error
		expectError bool
	}{
		{
			name: "adjust book meta",
			books: []bookdomain.Book{
				{
					ID:   "book",
					Name: "English Book 2",
				},
			},
			BookID:      "book",
			bookName:    "English Book",
			description: "Book for english lessons",
			grades:      []int{2},
			err:         nil,
			expectError: false,
		},
		{
			name: "adjust book meta name not set",
			books: []bookdomain.Book{
				{
					ID:   "book",
					Name: "English Book 2",
				},
			},
			BookID:      "book",
			bookName:    "",
			description: "Book for english lessons",
			grades:      []int{2},
			err:         bookdomain.ErrBookNameNotSet,
			expectError: true,
		},
		{
			name:        "adjust book meta book not found",
			books:       []bookdomain.Book{},
			BookID:      "book",
			bookName:    "English Book",
			description: "Book for english lessons",
			grades:      []int{2},
			err:         bookdomain.ErrBookWithIDNotFound("book"),
			expectError: true,
		},
		{
			name: "adjust book meta book with name Exists",
			books: []bookdomain.Book{
				{
					ID:   "book",
					Name: "English Book",
				},
			},
			BookID:      "book",
			bookName:    "English Book",
			description: "Book for english lessons",
			grades:      []int{2},
			err:         bookdomain.ErrBookWithNameAlreadyExists("English Book"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolBookAggregate(test.books)
			err := aggregate.AdjustBookMeta(test.BookID, test.bookName, test.description, test.grades)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, bookdomain.BookMetaAdjusted, event.EventType())
		})
	}
}

func TestIncreaseBookPrice(t *testing.T) {
	tests := []struct {
		name        string
		books       []bookdomain.Book
		BookID      string
		price       float64
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "increase book price",
			books: []bookdomain.Book{
				{
					ID:    "book",
					Price: 29.99,
				},
			},
			BookID:      "book",
			price:       49.99,
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name: "increase book price reason not specified",
			books: []bookdomain.Book{
				{
					ID:    "book",
					Price: 29.99,
				},
			},
			BookID:      "book",
			price:       49.99,
			reason:      "",
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
		{
			name: "increase book price, price less then zero",
			books: []bookdomain.Book{
				{
					ID:    "book",
					Price: 29.99,
				},
			},
			BookID:      "book",
			price:       -29.99,
			reason:      "test",
			err:         bookdomain.ErrBookPriceLessThanZero,
			expectError: true,
		},
		{
			name:        "increase book price, book not found",
			books:       []bookdomain.Book{},
			BookID:      "book",
			price:       39.99,
			reason:      "test",
			err:         bookdomain.ErrBookWithIDNotFound("book"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolBookAggregate(test.books)
			err := aggregate.IncreaseBookPrice(test.BookID, test.price, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, bookdomain.BookPriceIncreased, event.EventType())
		})
	}
}

func TestDecreaseBookPrice(t *testing.T) {
	tests := []struct {
		name        string
		books       []bookdomain.Book
		BookID      string
		price       float64
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "decrease book price",
			books: []bookdomain.Book{
				{
					ID:    "book",
					Price: 29.99,
				},
			},
			BookID:      "book",
			price:       49.99,
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name: "decrease book price reason not specified",
			books: []bookdomain.Book{
				{
					ID:    "book",
					Price: 29.99,
				},
			},
			BookID:      "book",
			price:       19.99,
			reason:      "",
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
		{
			name: "decrease book price, price less then zero",
			books: []bookdomain.Book{
				{
					ID:    "book",
					Price: 29.99,
				},
			},
			BookID:      "book",
			price:       -19.99,
			reason:      "test",
			err:         bookdomain.ErrBookPriceLessThanZero,
			expectError: true,
		},
		{
			name:        "decrease book price, book not found",
			books:       []bookdomain.Book{},
			BookID:      "book",
			price:       19.99,
			reason:      "test",
			err:         bookdomain.ErrBookWithIDNotFound("book"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolBookAggregate(test.books)
			err := aggregate.DecreaseBookPrice(test.BookID, test.price, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, bookdomain.BookPriceDecreased, event.EventType())
		})
	}
}
