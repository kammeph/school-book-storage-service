package common

import (
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Serializer interface {
	MarshalEvent(event common.Event) (Record, error)
	UnmarshalEvent(record Record) (common.Event, error)
}
