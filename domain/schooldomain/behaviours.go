package schooldomain

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

func (a *SchoolAggregate) AddSchool(name string) (string, error) {
	if name == "" {
		return "", ErrSchoolNameNotSet
	}
	for _, school := range a.Schools {
		if school.Name == name {
			return "", ErrSchoolAlreadyExists(name)
		}
	}
	schoolID := uuid.NewString()
	event, err := NewSchoolAdded(a, schoolID, name)
	if err != nil {
		return "", err
	}
	if err := a.Apply(event); err != nil {
		return "", err
	}
	return schoolID, nil
}

func (a *SchoolAggregate) DeactivateSchool(schoolID string, reason string) error {
	if !fp.Some(a.Schools, func(s School) bool { return s.ID == schoolID }) {
		return ErrSchoolWithIDNotFound(schoolID)
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	event, err := NewSchoolDeactivated(a, schoolID, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}

func (a *SchoolAggregate) RenameSchool(schoolID string, name string, reason string) error {
	if name == "" {
		return ErrSchoolNameNotSet
	}
	if reason == "" {
		return domain.ErrReasonNotSpecified
	}
	schoolsWithName := fp.Filter(a.Schools, func(s School) bool { return s.Name == name })
	if len(schoolsWithName) > 0 {
		for _, s := range schoolsWithName {
			if s.Name == name {
				return ErrSchoolAlreadyExists(s.Name)
			}
		}
	}
	event, err := NewSchoolRenamed(a, schoolID, name, reason)
	if err != nil {
		return err
	}
	if err := a.Apply(event); err != nil {
		return err
	}
	return nil
}
