package classdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain/classdomain"
	"github.com/stretchr/testify/assert"
)

func TestNewClass(t *testing.T) {
	id := "class"
	grade := 1
	letter := "A"
	pupils := 15
	from := time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC)
	timeStamp := time.Now()
	class := classdomain.NewClass(id, grade, letter, pupils, from, to, timeStamp)
	assert.Equal(t, id, class.ID)
	assert.Equal(t, grade, class.Grade)
	assert.Equal(t, letter, class.Letter)
	assert.Equal(t, pupils, class.NumberOfPupils)
	assert.EqualValues(t, from, class.DateFrom)
	assert.Equal(t, to, class.DateTo)
	assert.Equal(t, timeStamp, class.CreatedAt)
}
