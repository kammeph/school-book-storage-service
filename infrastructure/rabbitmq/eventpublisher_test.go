package rabbitmq_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/rabbitmq"
	"github.com/stretchr/testify/assert"
)

type EntityEvent struct {
	domain.EventModel
}

type EntityEvenHandler struct{}

func (h EntityEvenHandler) Handle(ctx context.Context, eventData []byte) {
	fmt.Printf("%v", eventData)
}

func TestNewEventPublisher(t *testing.T) {
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "create new event publisher",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "new event publisher channel error",
			connection:   NewMockConnection(true, false, false, false, false, false, false),
			err:          errChannel,
			exspectError: true,
		},
		{
			name:         "new event publisher exchange declare error",
			connection:   NewMockConnection(false, true, false, false, false, false, false),
			err:          errExchangeDeclare,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := rabbitmq.NewRabbitEventPublisher(test.connection, "test")
			if test.exspectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestPublish(t *testing.T) {
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "pulish message",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "publish message with error",
			connection:   NewMockConnection(false, false, false, false, false, true, false),
			err:          errPublish,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			broker, _ := rabbitmq.NewRabbitEventPublisher(test.connection, "test")
			event := &EntityEvent{
				EventModel: domain.EventModel{
					ID:      uuid.NewString(),
					Version: 1,
					At:      time.Now(),
				},
			}
			err := broker.Publish(context.Background(), []domain.Event{event})
			if test.exspectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
