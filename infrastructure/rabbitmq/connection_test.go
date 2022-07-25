package rabbitmq_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			err:          errChannel,
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
			err:          errClose,
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
