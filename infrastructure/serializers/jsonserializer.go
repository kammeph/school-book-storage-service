package serializers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
)

type JSONSerializer struct {
	eventTypes map[string]reflect.Type
}

type jsonEvent struct {
	Type string `json:"t"`
	Data string `json:"d"`
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{
		eventTypes: map[string]reflect.Type{},
	}
}

func NewJSONSerializerWithEvents(events ...domain.Event) *JSONSerializer {
	serializer := NewJSONSerializer()
	serializer.Bind(events...)
	return serializer
}

func (s *JSONSerializer) Bind(events ...domain.Event) {
	for _, event := range events {
		eventType, t := domain.EventType(event)
		s.eventTypes[eventType] = t
	}
}

func (s *JSONSerializer) MarshalEvent(event domain.Event) (common.Record, error) {
	eventType, _ := domain.EventType(event)
	data, err := json.Marshal(event)
	if err != nil {
		return common.Record{}, err
	}
	data, err = json.Marshal(jsonEvent{Type: eventType, Data: string(data)})
	if err != nil {
		return common.Record{}, nil
	}
	return common.Record{Version: event.EventVersion(), Data: string(data)}, nil
}

func (s *JSONSerializer) UnmarshalEvent(record common.Record) (domain.Event, error) {
	jEvent := jsonEvent{}
	err := json.Unmarshal([]byte(record.Data), &jEvent)
	if err != nil {
		return nil, err
	}
	t, ok := s.eventTypes[jEvent.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown event type: %s", jEvent.Type)
	}

	event := reflect.New(t).Interface()
	err = json.Unmarshal([]byte(jEvent.Data), event)
	if err != nil {
		return nil, err
	}
	return event.(domain.Event), nil
}
