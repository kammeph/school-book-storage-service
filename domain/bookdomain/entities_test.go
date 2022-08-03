package bookdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain/bookdomain"
	"github.com/stretchr/testify/assert"
)

func TestNewBook(t *testing.T) {
	id := "masterbook"
	isbn := "12345"
	name := "English Book"
	description := "book for english lessons"
	price := 25.90
	grades := []int{1, 2}
	timestamp := time.Now()
	book := bookdomain.NewBook(id, isbn, name, description, price, grades, timestamp)
	assert.Equal(t, id, book.ID)
	assert.Equal(t, isbn, book.Isbn)
	assert.Equal(t, name, book.Name)
	assert.Equal(t, description, book.Description)
	assert.Equal(t, price, book.Price)
	assert.Len(t, book.Grades, 2)
	assert.EqualValues(t, grades, book.Grades)
	assert.Equal(t, timestamp, book.CreatedAt)
	assert.Zero(t, book.UpdatedAt)
}
