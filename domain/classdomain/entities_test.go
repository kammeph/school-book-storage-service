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
	from := 2022
	to := 2023
	timeStamp := time.Now()
	class := classdomain.NewClass(id, grade, letter, from, to, timeStamp)
	assert.Equal(t, id, class.ID)
	assert.Equal(t, grade, class.Grade)
	assert.Equal(t, letter, class.Letter)
	assert.Equal(t, from, class.YearFrom)
	assert.Equal(t, to, class.YearTo)
	assert.Equal(t, timeStamp, class.CreatedAt)
}
