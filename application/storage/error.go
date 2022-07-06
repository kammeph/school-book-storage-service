package storage

import (
	"fmt"

	"github.com/kammeph/school-book-storage-service/domain/common"
)

func IncorrectAggregateTypeError(aggregate common.Aggregate) error {
	return fmt.Errorf("Incorrect type for aggregate: %T", aggregate)
}
