package domain_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/stretchr/testify/assert"
)

func TestSchoolEventCreation(t *testing.T) {
	tests := []struct {
		name       string
		schoolID   string
		schoolName string
		reason     string
		eventType  string
	}{
		{
			name:       "school added",
			schoolID:   "school",
			schoolName: "High School",
			eventType:  domain.SchoolAdded,
		},
		{
			name:      "storage deactivated",
			schoolID:  "school",
			reason:    "test",
			eventType: domain.SchoolDeactivated,
		},
		{
			name:       "school renamed",
			schoolID:   "school",
			schoolName: "High School",
			reason:     "test",
			eventType:  domain.SchoolRenamed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := domain.NewSchoolAggregate()
			var event domain.Event
			var err error
			switch test.eventType {
			case domain.SchoolAdded:
				event, err = domain.NewSchoolAdded(
					aggregate,
					test.schoolID,
					test.schoolName,
				)
			case domain.SchoolDeactivated:
				event, err = domain.NewSchoolDeactivated(
					aggregate,
					test.schoolID,
					test.reason,
				)
			case domain.SchoolRenamed:
				event, err = domain.NewSchoolRenamed(
					aggregate,
					test.schoolID,
					test.schoolName,
					test.reason,
				)
			default:
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, event)
			assert.Equal(t, aggregate.AggregateID(), event.AggregateID())
			assert.Equal(t, aggregate.AggregateVersion()+1, event.EventVersion())
			assert.NotZero(t, event.EventAt())
			assert.Equal(t, test.eventType, event.EventType())
			assert.NotEqual(t, "", event.EventData())
		})
	}
}
