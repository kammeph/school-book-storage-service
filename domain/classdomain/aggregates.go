package classdomain

import "github.com/kammeph/school-book-storage-service/domain"

type SchoolClassAggregate struct {
	*domain.AggregateModel
	Classes []Class
}

func (a *SchoolClassAggregate) On(event domain.Event) error {
	return nil
}
