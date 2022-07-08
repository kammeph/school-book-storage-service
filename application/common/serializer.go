package common

import (
	"github.com/kammeph/school-book-storage-service/domain/common"
)

type Serializer interface {
	Bind(events ...common.Event)
	MarshalEvent(event common.Event) (Record, error)
	UnmarshalEvent(record Record) (common.Event, error)
}
