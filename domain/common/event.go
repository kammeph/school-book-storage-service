package common

import (
	"encoding/json"
	"time"
)

// History represents
type History []Event

// Len implements sort.Interface
func (h History) Len() int {
	return len(h)
}

// Swap implements sort.Interface
func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Less implements sort.Interface
func (h History) Less(i, j int) bool {
	return h[i].EventVersion() < h[j].EventVersion()
}

type Event interface {
	AggregateID() string
	EventVersion() int
	EventAt() time.Time
	EventType() string
	EventData() string
	GetJsonData(data interface{}) error
	SetJsonData(data interface{}) error
}

type EventModel struct {
	ID      string
	Version int
	At      time.Time
	Type    string
	Data    string
}

func NewEvent(aggregate Aggregate, eventType string) Event {
	return &EventModel{
		ID:      aggregate.AggregateID(),
		Version: aggregate.AggregateVersion() + 1,
		At:      time.Now(),
		Type:    eventType,
	}
}

func (m EventModel) AggregateID() string {
	return m.ID
}

func (m EventModel) EventVersion() int {
	return m.Version
}

func (m EventModel) EventAt() time.Time {
	return m.At
}

func (m EventModel) EventType() string {
	return m.Type
}

func (m EventModel) EventData() string {
	return m.Data
}

func (m EventModel) GetJsonData(data interface{}) error {
	return json.Unmarshal([]byte(m.Data), data)
}

func (m *EventModel) SetJsonData(data interface{}) error {
	eventData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	m.Data = string(eventData)
	return nil
}
