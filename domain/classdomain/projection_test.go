package classdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain/classdomain"
	"github.com/stretchr/testify/assert"
)

func TestNewClassWithBooks(t *testing.T) {
	schoolID := "school"
	classID := "class"
	grade := 1
	letter := "A"
	pupils := 15
	from := time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC)
	class := classdomain.NewClassWithBooks(schoolID, classID, grade, letter, pupils, from, to)
	assert.Equal(t, schoolID, class.SchoolID)
	assert.Equal(t, classID, class.ClassID)
	assert.Equal(t, grade, class.Grade)
	assert.Equal(t, letter, class.Letter)
	assert.Equal(t, pupils, class.NumberOfPupils)
	assert.EqualValues(t, from, class.DateFrom)
	assert.Equal(t, to, class.DateTo)
	assert.NotNil(t, class.Books)
	assert.Len(t, class.Books, 0)
}
