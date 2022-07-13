package common

import domain "github.com/kammeph/school-book-storage-service/domain/common"

type MessageBroker interface {
	Publish(event domain.Event) error
	Subscribe(event domain.Event, handler EventHandler) error
}
