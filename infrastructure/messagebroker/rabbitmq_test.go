package messagebroker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/messagebroker"
	"github.com/stretchr/testify/assert"
)

type EntityEvent struct {
	domain.EventModel
}

type EntityEvenHandler struct{}

func (h EntityEvenHandler) Handle(ctx context.Context, eventData []byte) {
	fmt.Printf("%v", eventData)
}

func TestNewRabbitMQ(t *testing.T) {
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "create new rabbit mq",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "new rabbits mq channel error",
			connection:   NewMockConnection(true, false, false, false, false, false, false),
			err:          channelError,
			exspectError: true,
		},
		{
			name:         "new rabbits mq exchange declare error",
			connection:   NewMockConnection(false, true, false, false, false, false, false),
			err:          exchangeDeclareError,
			exspectError: true,
		},
		{
			name:         "new rabbits mq queue declare error",
			connection:   NewMockConnection(false, false, true, false, false, false, false),
			err:          queueDeclareError,
			exspectError: true,
		},
		{
			name:         "new rabbits mq queue bind error",
			connection:   NewMockConnection(false, false, false, true, false, false, false),
			err:          queueBindError,
			exspectError: true,
		},
		{
			name:         "new rabbits mq consume error",
			connection:   NewMockConnection(false, false, false, false, true, false, false),
			err:          consumeError,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := messagebroker.NewRabbitMQ(test.connection, "test")
			if test.exspectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestSubscribe(t *testing.T) {
	conn := NewMockConnection(false, false, false, false, false, false, false)
	broker, err := messagebroker.NewRabbitMQ(conn, "test")
	assert.Nil(t, err)
	assert.NotNil(t, broker)
	broker.Subscribe(EntityEvent{}, EntityEvenHandler{})
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
			err:          publishError,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			broker, _ := messagebroker.NewRabbitMQ(test.connection, "test")
			event := EntityEvent{
				EventModel: domain.EventModel{
					ID:      uuid.NewString(),
					Version: 1,
					At:      time.Now(),
				},
			}
			err := broker.Publish(event)
			if test.exspectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetChannel(t *testing.T) {
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "get channel",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "get channel with error",
			connection:   NewMockConnection(true, false, false, false, false, true, false),
			err:          channelError,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.connection.Channel()
			if test.exspectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestConnectionClose(t *testing.T) {
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "get channel",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "get channel with error",
			connection:   NewMockConnection(false, false, false, false, false, false, true),
			err:          closeError,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.connection.Close()
			if test.exspectError {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
