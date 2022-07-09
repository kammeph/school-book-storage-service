package serializers_test

import (
	"testing"

	"github.com/google/uuid"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/serializers"
	"github.com/stretchr/testify/assert"
)

type EntityNameSet struct {
	domain.EventModel
	Name string
}

func TestJSONSerializer(t *testing.T) {
	event := EntityNameSet{
		EventModel: domain.EventModel{
			ID: uuid.New().String(), Version: 5,
		}, Name: "Test",
	}
	serializer := serializers.NewJSONSerializerWithEvents(event)
	record, err := serializer.MarshalEvent(event)
	assert.Nil(t, err)
	assert.NotNil(t, record)

	v, err := serializer.UnmarshalEvent(record)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	found, ok := v.(*EntityNameSet)
	assert.True(t, ok)
	assert.Equal(t, &event, found)
}
