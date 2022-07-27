package classdomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/classdomain"
	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/stretchr/testify/assert"
)

type UnknownEvent struct {
	domain.EventModel
}

func TestSchoolClassAggregateOn(t *testing.T) {
	tests := []struct {
		name              string
		eventType         string
		eventVersion      int
		eventAt           time.Time
		schoolID          string
		classID           string
		grade             int
		letter            string
		numberOfPupils    int
		dateFrom          time.Time
		dateTo            time.Time
		reason            string
		err               error
		expectError       bool
		addDefaultClasses bool
	}{
		{
			name:              "class created",
			eventType:         classdomain.ClassCreated,
			eventVersion:      4,
			eventAt:           time.Now(),
			schoolID:          "school",
			classID:           "2A",
			grade:             2,
			letter:            "A",
			numberOfPupils:    16,
			dateFrom:          time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:            time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			err:               nil,
			expectError:       false,
			addDefaultClasses: false,
		},
		{
			name:              "class created error",
			eventType:         classdomain.ClassCreated,
			eventVersion:      9,
			eventAt:           time.Now(),
			schoolID:          "school",
			classID:           "2A",
			grade:             2,
			letter:            "A",
			numberOfPupils:    16,
			dateFrom:          time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			dateTo:            time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
			err:               classdomain.ErrApplyEventClassAlreadyExists(classdomain.ClassCreated, "2A"),
			expectError:       true,
			addDefaultClasses: true,
		},
		{
			name:              "number of pupils increased",
			eventType:         classdomain.NumberOfPupilsIncreased,
			eventVersion:      7,
			eventAt:           time.Now(),
			classID:           "1C",
			numberOfPupils:    16,
			err:               nil,
			expectError:       false,
			addDefaultClasses: true,
		},
		{
			name:              "number of pupils increased error",
			eventType:         classdomain.NumberOfPupilsIncreased,
			eventVersion:      9,
			eventAt:           time.Now(),
			classID:           "4B",
			numberOfPupils:    16,
			err:               classdomain.ErrApplyEventClassNotFound(classdomain.NumberOfPupilsIncreased, "4B"),
			expectError:       true,
			addDefaultClasses: false,
		},
		{
			name:              "number of pupils decreased",
			eventType:         classdomain.NumberOfPupilsDecreased,
			eventVersion:      21,
			eventAt:           time.Now(),
			classID:           "3D",
			numberOfPupils:    21,
			err:               nil,
			expectError:       false,
			addDefaultClasses: true,
		},
		{
			name:              "number of pupils decreased error",
			eventType:         classdomain.NumberOfPupilsDecreased,
			eventVersion:      52,
			eventAt:           time.Now(),
			classID:           "5E",
			numberOfPupils:    14,
			err:               classdomain.ErrApplyEventClassNotFound(classdomain.NumberOfPupilsDecreased, "5E"),
			expectError:       true,
			addDefaultClasses: false,
		},
		{
			name:              "unknown",
			eventType:         "unknown",
			eventVersion:      52,
			eventAt:           time.Now(),
			err:               domain.ErrUnknownEvent(&UnknownEvent{}),
			expectError:       true,
			addDefaultClasses: false,
		},
	}
	for _, test := range tests {
		var event domain.Event
		t.Run(test.name, func(t *testing.T) {
			switch test.eventType {
			case classdomain.ClassCreated:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := classdomain.ClassCreatedEvent{
					test.schoolID,
					test.classID,
					test.grade,
					test.letter,
					test.numberOfPupils,
					test.dateFrom,
					test.dateTo,
				}
				event.SetJsonData(eventData)
			case classdomain.NumberOfPupilsIncreased:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := classdomain.NumberOfPupilsIncreasedEvent{
					test.classID,
					test.numberOfPupils,
					test.reason,
				}
				event.SetJsonData(eventData)
			case classdomain.NumberOfPupilsDecreased:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := classdomain.NumberOfPupilsDecreasedEvent{
					test.classID,
					test.numberOfPupils,
					test.reason,
				}
				event.SetJsonData(eventData)
			default:
				event = &UnknownEvent{}
			}
			aggregate := classdomain.NewSchoolClassAggregate()
			if test.addDefaultClasses {
				aggregate.Classes = append(aggregate.Classes, classdomain.Class{
					ID:             test.classID,
					NumberOfPupils: 25,
				})
			}
			err := aggregate.On(event)
			if test.expectError {
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.eventVersion, aggregate.Version)
			class := fp.Find(aggregate.Classes, func(c classdomain.Class) bool { return c.ID == test.classID })
			assert.NotNil(t, class)
			if test.eventType == classdomain.ClassCreated {
				assert.Equal(t, test.eventVersion, aggregate.Version)
				assert.Equal(t, test.grade, class.Grade)
				assert.Equal(t, test.letter, class.Letter)
				assert.Equal(t, test.numberOfPupils, class.NumberOfPupils)
				assert.Equal(t, test.dateFrom, class.DateFrom)
				assert.Equal(t, test.dateTo, class.DateTo)
				assert.Equal(t, test.eventAt, class.CreatedAt)
				return
			}
			if test.eventType == classdomain.NumberOfPupilsIncreased {
				assert.Equal(t, test.eventVersion, aggregate.Version)
				assert.Equal(t, 25+test.numberOfPupils, class.NumberOfPupils)
				assert.Equal(t, test.eventAt, class.UpdatedAt)
			}
			if test.eventType == classdomain.NumberOfPupilsDecreased {
				assert.Equal(t, test.eventVersion, aggregate.Version)
				assert.Equal(t, 25-test.numberOfPupils, class.NumberOfPupils)
				assert.Equal(t, test.eventAt, class.UpdatedAt)
			}
		})
	}
}
