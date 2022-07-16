package common

type on func(event Event) error

type AggregateModel struct {
	ID      string
	Version int
	Events  []Event
	on      on
}

func NewAggregateModel(on on) AggregateModel {
	return AggregateModel{on: on, Events: []Event{}}
}

func (a AggregateModel) AggregateID() string {
	return a.ID
}

func (a *AggregateModel) SetAggregateID(id string) {
	a.ID = id
}

func (a AggregateModel) AggregateVersion() int {
	return a.Version
}

func (a AggregateModel) DomainEvents() []Event {
	return a.Events
}

func (a AggregateModel) Load(events []Event) error {
	for _, event := range events {
		err := a.on(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AggregateModel) Apply(event Event) error {
	err := a.on(event)
	if err != nil {
		return err
	}
	a.Events = append(a.Events, event)
	return nil
}

type Aggregate interface {
	AggregateID() string
	SetAggregateID(string)
	AggregateVersion() int
	DomainEvents() []Event
	Load(events []Event) error
	Apply(event Event) error
	On(event Event) error
}
