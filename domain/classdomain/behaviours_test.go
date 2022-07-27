package classdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/classdomain"
	"github.com/stretchr/testify/assert"
)

func initSchoolClassAggregate(classes []classdomain.Class) *classdomain.SchoolClassAggregate {
	aggregate := classdomain.NewSchoolClassAggregate()
	aggregate.Classes = classes
	return aggregate
}

func TestAddSchool(t *testing.T) {
	tests := []struct {
		name        string
		classes     []classdomain.Class
		grade       int
		letter      string
		number      int
		dateFrom    time.Time
		dateTo      time.Time
		err         error
		expectError bool
	}{
		{
			name:        "create class",
			classes:     []classdomain.Class{},
			grade:       2,
			letter:      "A",
			number:      16,
			dateFrom:    time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         nil,
			expectError: false,
		},
		{
			name:        "create class invalid grade",
			classes:     []classdomain.Class{},
			grade:       0,
			letter:      "A",
			number:      16,
			dateFrom:    time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         classdomain.ErrGradeGreaterZero,
			expectError: true,
		},
		{
			name:        "create class letter not set",
			classes:     []classdomain.Class{},
			grade:       4,
			letter:      "",
			number:      16,
			dateFrom:    time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         classdomain.ErrLetterNotSet,
			expectError: true,
		},
		{
			name:        "create class letter more than one place",
			classes:     []classdomain.Class{},
			grade:       4,
			letter:      "AB",
			number:      16,
			dateFrom:    time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         classdomain.ErrLetterToLong,
			expectError: true,
		},
		{
			name:        "create class pupils less than one",
			classes:     []classdomain.Class{},
			grade:       4,
			letter:      "A",
			number:      0,
			dateFrom:    time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         classdomain.ErrNumberOfPupilsGreaterZero,
			expectError: true,
		},
		{
			name:        "create class date from greater than date to",
			classes:     []classdomain.Class{},
			grade:       4,
			letter:      "A",
			number:      9,
			dateFrom:    time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         classdomain.ErrInvalidDates,
			expectError: true,
		},
		{
			name: "create class date from greater than date to",
			classes: []classdomain.Class{
				{
					ID:             "4A",
					Grade:          4,
					Letter:         "A",
					NumberOfPupils: 9,
					DateFrom:       time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
					DateTo:         time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			grade:       4,
			letter:      "A",
			number:      9,
			dateFrom:    time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:      time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			err:         classdomain.ErrClassAlreadyExists(4, "A", 2022, 2023),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolClassAggregate(test.classes)
			classID, err := aggregate.CreateClass(test.grade, test.letter, test.number, test.dateFrom, test.dateTo)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			assert.NotEqual(t, "", classID)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, classdomain.ClassCreated, event.EventType())
		})
	}
}

func TestIncreaseNumberOfPupils(t *testing.T) {
	tests := []struct {
		name        string
		classes     []classdomain.Class
		classID     string
		number      int
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "increase number of pupils",
			classes: []classdomain.Class{
				{
					ID:             "4A",
					Grade:          4,
					Letter:         "A",
					NumberOfPupils: 9,
					DateFrom:       time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
					DateTo:         time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			classID:     "4A",
			number:      16,
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name:        "increase number of pupils with zero",
			classes:     []classdomain.Class{},
			classID:     "4A",
			number:      0,
			err:         classdomain.ErrIncreasePupilsGreaterZero,
			expectError: true,
		},
		{
			name:        "increase number of pupils with no reason",
			classes:     []classdomain.Class{},
			classID:     "4A",
			number:      16,
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
		{
			name:        "increase number of pupils class not exists",
			classes:     []classdomain.Class{},
			classID:     "4A",
			number:      16,
			reason:      "test",
			err:         classdomain.ErrClassWithIDNotFound("4A"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolClassAggregate(test.classes)
			err := aggregate.IncreaseNumberOfPupils(test.classID, test.number, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, classdomain.NumberOfPupilsIncreased, event.EventType())
		})
	}
}

func TestDecreaseNumberOfPupils(t *testing.T) {
	tests := []struct {
		name        string
		classes     []classdomain.Class
		classID     string
		number      int
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "decrease number of pupils",
			classes: []classdomain.Class{
				{
					ID:             "4A",
					Grade:          4,
					Letter:         "A",
					NumberOfPupils: 9,
					DateFrom:       time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
					DateTo:         time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			classID:     "4A",
			number:      16,
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name:        "decrease number of pupils with zero",
			classes:     []classdomain.Class{},
			classID:     "4A",
			number:      0,
			err:         classdomain.ErrDecreasePupilsGreaterZero,
			expectError: true,
		},
		{
			name:        "decrease number of pupils with no reason",
			classes:     []classdomain.Class{},
			classID:     "4A",
			number:      16,
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
		{
			name:        "decrease number of pupils class not exists",
			classes:     []classdomain.Class{},
			classID:     "4A",
			number:      16,
			reason:      "test",
			err:         classdomain.ErrClassWithIDNotFound("4A"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolClassAggregate(test.classes)
			err := aggregate.DecreaseNumberOfPupils(test.classID, test.number, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, classdomain.NumberOfPupilsDecreased, event.EventType())
		})
	}
}
