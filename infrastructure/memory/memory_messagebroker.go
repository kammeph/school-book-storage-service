package memory

import (
	"context"
	"encoding/json"

	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
)

type MemoryMessageBroker struct {
	eventHandlers []common.EventHandler
}

func NewMemoryMessageBroker() *MemoryMessageBroker {
	return &MemoryMessageBroker{}
}

func (m *MemoryMessageBroker) Publish(ctx context.Context, events []domain.Event) error {
	for _, event := range events {
		eventBytes, err := json.Marshal(event)
		if err != nil {
			return err
		}
		for idx := range m.eventHandlers {
			m.eventHandlers[idx].Handle(ctx, []byte(eventBytes))
		}
	}
	return nil
}

func (m *MemoryMessageBroker) Subscribe(exchange string, handler common.EventHandler) error {
	m.eventHandlers = append(m.eventHandlers, handler)
	return nil
}
