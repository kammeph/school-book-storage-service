package schooldomain_test

import (
	"testing"
	"time"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
	"github.com/stretchr/testify/assert"
)

type UnknownEvent struct {
	domain.EventModel
}

func TestSchoolOn(t *testing.T) {
	tests := []struct {
		name              string
		eventVersion      int
		eventAt           time.Time
		schoolID          string
		schoolName        string
		reason            string
		eventType         string
		err               error
		expectError       bool
		addDefaultStorage bool
		operation         string
	}{
		{
			name:              "on school added",
			eventVersion:      1,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         schooldomain.SchoolAdded,
			expectError:       false,
			addDefaultStorage: false,
			operation:         "add",
		},
		{
			name:              "try add school twice",
			eventVersion:      1,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         schooldomain.SchoolAdded,
			err:               schooldomain.ErrSchoolWithIdAlreadyExists("school"),
			expectError:       true,
			addDefaultStorage: true,
		},
		{
			name:              "on school deactivate",
			eventVersion:      7,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         schooldomain.SchoolDeactivated,
			err:               schooldomain.ErrSchoolWithIDNotFound("school"),
			expectError:       false,
			addDefaultStorage: true,
			operation:         "deactivate",
		},
		{
			name:              "deactivate non existing school",
			eventVersion:      34,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         schooldomain.SchoolDeactivated,
			err:               schooldomain.ErrSchoolWithIDNotFound("school"),
			expectError:       true,
			addDefaultStorage: false,
			operation:         "deactivate",
		},
		{
			name:              "on school renamed",
			eventVersion:      5,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         schooldomain.SchoolRenamed,
			err:               nil,
			expectError:       false,
			addDefaultStorage: true,
			operation:         "update",
		},
		{
			name:              "rename non existing school",
			eventVersion:      3,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         schooldomain.SchoolRenamed,
			err:               schooldomain.ErrSchoolWithIDNotFound("school"),
			expectError:       true,
			addDefaultStorage: false,
		},
		{
			name:              "unknown event",
			eventVersion:      9,
			eventAt:           time.Now(),
			schoolID:          "school",
			schoolName:        "High School",
			reason:            "test",
			eventType:         "UnknownEvent",
			err:               domain.ErrUnknownEvent(&UnknownEvent{}),
			expectError:       true,
			addDefaultStorage: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event domain.Event
			switch test.eventType {
			case schooldomain.SchoolAdded:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := schooldomain.SchoolAddedEvent{test.schoolID, test.schoolName}
				event.SetJsonData(eventData)
			case schooldomain.SchoolDeactivated:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := schooldomain.SchoolDeactivatedEvent{test.schoolID, test.reason}
				event.SetJsonData(eventData)
			case schooldomain.SchoolRenamed:
				event = &domain.EventModel{
					Version: test.eventVersion,
					At:      test.eventAt,
					Type:    test.eventType,
				}
				eventData := schooldomain.SchoolRenamedEvent{test.schoolID, test.schoolName, test.reason}
				event.SetJsonData(eventData)
			default:
				event = &UnknownEvent{}
			}
			aggregate := schooldomain.NewSchoolAggregate()
			if test.addDefaultStorage {
				aggregate.Schools = append(aggregate.Schools, schooldomain.School{
					ID:     test.schoolID,
					Name:   "High School",
					Active: true,
				})
			}
			err := aggregate.On(event)
			if test.expectError {
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.eventVersion, aggregate.Version)
			school, idx, err := aggregate.GetSchoolByID(test.schoolID)
			if test.operation == "add" {
				assert.Equal(t, test.eventAt, school.CreatedAt)
			}
			if test.operation == "deactivate" {
				assert.Equal(t, test.eventAt, school.UpdatedAt)
				assert.False(t, school.Active)
			}
			if test.operation == "update" {
				assert.Equal(t, test.eventAt, school.UpdatedAt)
			}
			assert.NoError(t, err)
			assert.Greater(t, idx, -1)
			assert.NotNil(t, school)
			assert.Equal(t, test.schoolName, school.Name)
		})
	}
}
