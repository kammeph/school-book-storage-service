package classdomain

type ClassCreatedEvent struct {
	SchoolID string `json:"schoolId"`
	ClassID  string `json:"classId"`
	Grade    int    `json:"grade"`
	Letter   string `json:"letter"`
	YearFrom int    `json:"yearFrom"`
	YearTo   int    `json:"yearTo"`
}

func NewClassCreatedt(aggregate *SchoolClassAggregate, classID string, grade int, letter string, yearFrom, yearTo int) ClassCreatedEvent {
	return ClassCreatedEvent{
		SchoolID: aggregate.AggregateID(),
		ClassID:  classID,
		Grade:    grade,
		Letter:   letter,
		YearFrom: yearFrom,
		YearTo:   yearTo,
	}
}
