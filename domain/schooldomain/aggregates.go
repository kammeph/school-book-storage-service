package schooldomain

import (
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/fp"
)

type SchoolAggregate struct {
	*domain.AggregateModel
	Schools []School
}

func NewSchoolAggregate() *SchoolAggregate {
	aggregate := &SchoolAggregate{
		Schools: []School{},
	}
	model := domain.NewAggregateModel(aggregate.On)
	aggregate.AggregateModel = &model
	return aggregate
}

func (a *SchoolAggregate) On(event domain.Event) error {
	switch event.EventType() {
	case SchoolAdded:
		return a.onSchoolAdded(event)
	case SchoolDeactivated:
		return a.onSchoolDeactivated(event)
	case SchoolRenamed:
		return a.onSchoolRenamed(event)
	default:
		return domain.ErrUnknownEvent(event)
	}
}

func (a *SchoolAggregate) onSchoolAdded(event domain.Event) error {
	eventData := SchoolAddedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	if fp.Some(a.Schools, func(s School) bool { return s.ID == eventData.SchoolID }) {
		return ErrSchoolWithIdAlreadyExists(eventData.SchoolID)
	}
	school := NewSchool(eventData.SchoolID, eventData.Name, event.EventAt())
	a.Version = event.EventVersion()
	a.Schools = append(a.Schools, school)
	return nil
}

func (a *SchoolAggregate) onSchoolDeactivated(event domain.Event) error {
	eventData := SchoolDeactivatedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	school := fp.Find(a.Schools, func(s School) bool { return s.ID == eventData.SchoolID })
	if school == nil {
		return ErrSchoolWithIDNotFound(eventData.SchoolID)
	}
	a.Version = event.EventVersion()
	school.UpdatedAt = event.EventAt()
	school.Active = false
	return nil
}

func (a *SchoolAggregate) onSchoolRenamed(event domain.Event) error {
	eventData := SchoolRenamedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		return err
	}
	school := fp.Find(a.Schools, func(s School) bool { return s.ID == eventData.SchoolID })
	if school == nil {
		return ErrSchoolWithIDNotFound(eventData.SchoolID)
	}
	a.Version = event.EventVersion()
	school.UpdatedAt = event.EventAt()
	school.Name = eventData.Name
	return nil
}

// func (a *SchoolAggregate) GetSchoolByID(id string) (*School, int, error) {
// 	for idx, s := range a.Schools {
// 		if s.ID == id {
// 			return &a.Schools[idx], idx, nil
// 		}
// 	}
// 	return nil, -1, ErrSchoolWithIDNotFound(id)
// }

// func (a *SchoolAggregate) RemoveSchoolByID(id string) error {
// 	_, idx, err := a.GetSchoolByID(id)
// 	if err != nil {
// 		return err
// 	}
// 	a.Schools = append(a.Schools[:idx], a.Schools[idx+1:]...)
// 	return nil
// }

// func (a SchoolAggregate) GetSchoolByName(name string) (*School, error) {
// 	schools := a.getSchoolsByName(name)
// 	if len(schools) == 0 {
// 		return nil, ErrSchoolWithNameNotFound(name)
// 	}
// 	if len(schools) > 1 {
// 		return nil, ErrMultipleSchoolsWithNameFound(name)
// 	}
// 	return &schools[0], nil
// }

// func (a SchoolAggregate) getSchoolsByName(name string) []School {
// 	schools := []School{}
// 	for _, storage := range a.Schools {
// 		if storage.Name == name {
// 			schools = append(schools, storage)
// 		}
// 	}
// 	return schools
// }
