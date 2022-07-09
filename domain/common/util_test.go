package common_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/stretchr/testify/assert"
)

type EntityNameSet struct {
	common.EventModel
	Name string
}

type Custom struct {
	common.EventModel
	Name string
}

func (c Custom) EventType() string {
	return "custom"
}

func TestEventType(t *testing.T) {
	event := EntityNameSet{
		EventModel: common.EventModel{
			ID: uuid.New().String(), Version: 5,
		}, Name: "Test",
	}
	eventTypeName, _ := common.EventType(event)
	assert.Equal(t, eventTypeName, "EntityNameSet")
}

func TestEventTypeCustomTyper(t *testing.T) {
	event := Custom{
		EventModel: common.EventModel{
			ID: uuid.New().String(), Version: 5,
		},
	}
	eventTypeName, _ := common.EventType(event)
	assert.Equal(t, eventTypeName, "custom")
}
