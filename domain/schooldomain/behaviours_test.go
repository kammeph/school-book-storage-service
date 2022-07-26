package schooldomain_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
	"github.com/stretchr/testify/assert"
)

func initSchoolAggregate(schools []schooldomain.School) *schooldomain.SchoolAggregate {
	aggregate := schooldomain.NewSchoolAggregate()
	aggregate.Schools = schools
	return aggregate
}

func TestAddSchool(t *testing.T) {
	tests := []struct {
		name        string
		schools     []schooldomain.School
		schoolname  string
		err         error
		expectError bool
	}{
		{
			name:        "add school",
			schools:     []schooldomain.School{},
			schoolname:  "High School",
			err:         nil,
			expectError: false,
		},
		{
			name:        "add school without name",
			schools:     []schooldomain.School{},
			schoolname:  "",
			err:         schooldomain.ErrSchoolNameNotSet,
			expectError: true,
		},
		{
			name: "school already exists",
			schools: []schooldomain.School{{
				ID:   "school",
				Name: "High School",
			}},
			schoolname:  "High School",
			err:         schooldomain.ErrSchoolAlreadyExists("High School"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolAggregate(test.schools)
			schoolID, err := aggregate.AddSchool(test.schoolname)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			assert.NotEqual(t, "", schoolID)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, schooldomain.SchoolAdded, event.EventType())
		})
	}
}

func TestDeactivateSchool(t *testing.T) {
	tests := []struct {
		name        string
		schools     []schooldomain.School
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "deactivate school",
			schools: []schooldomain.School{{
				ID:     "school",
				Name:   "High School",
				Active: true,
			}},
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name:        "error when deactivating not existing school",
			schools:     []schooldomain.School{},
			reason:      "test",
			err:         schooldomain.ErrSchoolWithIDNotFound("school"),
			expectError: true,
		},
		{
			name: "error when deactivating without a reason",
			schools: []schooldomain.School{{
				ID:     "school",
				Name:   "High School",
				Active: true,
			}},
			reason:      "",
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolAggregate(test.schools)
			err := aggregate.DeactivateSchool("school", test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			assert.Len(t, aggregate.Schools, 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, schooldomain.SchoolDeactivated, event.EventType())
		})
	}
}

func TestRenameSchool(t *testing.T) {
	tests := []struct {
		name        string
		schools     []schooldomain.School
		schoolName  string
		reason      string
		err         error
		expectError bool
	}{
		{
			name: "rename school",
			schools: []schooldomain.School{{
				ID:   "school",
				Name: "High School",
			}},
			schoolName:  "Mid School",
			reason:      "test",
			err:         nil,
			expectError: false,
		},
		{
			name: "storage with same name exists",
			schools: []schooldomain.School{{
				ID:   "school",
				Name: "High School",
			}},
			schoolName:  "High School",
			reason:      "test",
			err:         schooldomain.ErrSchoolAlreadyExists("High School"),
			expectError: true,
		},
		{
			name: "school name not set",
			schools: []schooldomain.School{{
				ID:   "school",
				Name: "storage",
			}},
			schoolName:  "",
			reason:      "test",
			err:         schooldomain.ErrSchoolNameNotSet,
			expectError: true,
		},
		{
			name: "reason not specified",
			schools: []schooldomain.School{{
				ID:   "school",
				Name: "High School",
			}},
			schoolName:  "Mid School",
			reason:      "",
			err:         domain.ErrReasonNotSpecified,
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aggregate := initSchoolAggregate(test.schools)
			err := aggregate.RenameSchool("school", test.schoolName, test.reason)
			if test.expectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, aggregate.DomainEvents(), 1)
			event := aggregate.DomainEvents()[0]
			assert.Equal(t, schooldomain.SchoolRenamed, event.EventType())
		})
	}
}
