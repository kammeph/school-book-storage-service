package schooldomain_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
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
			eventType:  schooldomain.SchoolAdded,
		},
		{
			name:      "storage deactivated",
			schoolID:  "school",
			reason:    "test",
			eventType: schooldomain.SchoolDeactivated,
		},
		{
			name:       "school renamed",
			schoolID:   "school",
			schoolName: "High School",
			reason:     "test",
			eventType:  schooldomain.SchoolRenamed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := schooldomain.NewSchoolAggregate()
			var event domain.Event
			var err error
			switch test.eventType {
			case schooldomain.SchoolAdded:
				event, err = schooldomain.NewSchoolAdded(
					aggregate,
					test.schoolID,
					test.schoolName,
				)
			case schooldomain.SchoolDeactivated:
				event, err = schooldomain.NewSchoolDeactivated(
					aggregate,
					test.schoolID,
					test.reason,
				)
			case schooldomain.SchoolRenamed:
				event, err = schooldomain.NewSchoolRenamed(
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
