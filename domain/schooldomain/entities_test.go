package schooldomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
	"github.com/stretchr/testify/assert"
)

func TestNewSchool(t *testing.T) {
	id := "NewSchool"
	name := "High School"
	timeStamp := time.Now()
	school := schooldomain.NewSchool(id, name, timeStamp)
	assert.NotZero(t, school)
	assert.Equal(t, id, school.ID)
	assert.Equal(t, name, school.Name)
	assert.True(t, school.Active)
	assert.Equal(t, timeStamp, school.CreatedAt)
	assert.Zero(t, school.UpdatedAt)
}
