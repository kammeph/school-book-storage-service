package application

import (
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain"
)

func ErrIncorrectAggregateType(aggregate domain.Aggregate) error {
	return fmt.Errorf("incorrect type for aggregate: %T", aggregate)
}
