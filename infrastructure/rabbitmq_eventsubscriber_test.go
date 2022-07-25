package infrastructure_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestNewEventSubscriber(t *testing.T) {
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "create new event subscriber",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "new event subscriber channel error",
			connection:   NewMockConnection(true, false, false, false, false, false, false),
			err:          errChannel,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := infrastructure.NewRabbitEventSubscriber(test.connection)
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
	tests := []struct {
		name         string
		connection   MockConntection
		err          error
		exspectError bool
	}{
		{
			name:         "subscribe",
			connection:   NewMockConnection(false, false, false, false, false, false, false),
			err:          nil,
			exspectError: false,
		},
		{
			name:         "subscribe queue declare error",
			connection:   NewMockConnection(false, false, true, false, false, false, false),
			err:          errQueueDeclare,
			exspectError: true,
		},
		{
			name:         "subscribe queue bind error",
			connection:   NewMockConnection(false, false, false, true, false, false, false),
			err:          errQueueBind,
			exspectError: true,
		},
		{
			name:         "subscribe consume error",
			connection:   NewMockConnection(false, false, false, false, true, false, false),
			err:          errConsume,
			exspectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			broker, err := infrastructure.NewRabbitEventSubscriber(test.connection)
			assert.Nil(t, err)
			assert.NotNil(t, broker)
			broker.Subscribe("test", EntityEvenHandler{})
		})
	}
}
