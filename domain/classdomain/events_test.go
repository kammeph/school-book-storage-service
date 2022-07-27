package classdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain/classdomain"
	"github.com/stretchr/testify/assert"
)

func TestNewClassCreated(t *testing.T) {
	classID := "3A"
	grade := 3
	letter := "A"
	number := 17
	dateFrom := time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC)
	aggregate := classdomain.NewSchoolClassAggregate()
	aggregate.Version = 4
	event, err := classdomain.NewClassCreated(aggregate, classID, grade, letter, number, dateFrom, dateTo)
	assert.NoError(t, err)
	assert.Equal(t, 5, event.EventVersion())
	assert.Equal(t, classdomain.ClassCreated, event.EventType())
	assert.NotZero(t, event.EventAt())
	assert.NotEqual(t, "", event.EventData())
}

func TestNewNumberOfPupilsIncreased(t *testing.T) {
	classID := "2B"
	number := 2
	reason := "test"
	aggregate := classdomain.NewSchoolClassAggregate()
	aggregate.Version = 7
	event, err := classdomain.NewNumberOfPupilsIncreased(aggregate, classID, number, reason)
	assert.NoError(t, err)
	assert.Equal(t, 8, event.EventVersion())
	assert.Equal(t, classdomain.NumberOfPupilsIncreased, event.EventType())
	assert.NotZero(t, event.EventAt())
	assert.NotEqual(t, "", event.EventData())
}

func TestNewNumberOfPupilsDecreased(t *testing.T) {
	classID := "2B"
	number := 2
	reason := "test"
	aggregate := classdomain.NewSchoolClassAggregate()
	aggregate.Version = 9
	event, err := classdomain.NewNumberOfPupilsDecreased(aggregate, classID, number, reason)
	assert.NoError(t, err)
	assert.Equal(t, 10, event.EventVersion())
	assert.Equal(t, classdomain.NumberOfPupilsDecreased, event.EventType())
	assert.NotZero(t, event.EventAt())
	assert.NotEqual(t, "", event.EventData())
}
