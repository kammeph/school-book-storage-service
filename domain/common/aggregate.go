package common

type AggregateModel struct {
	ID      string
	Version int
	Events  []Event
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

type Aggregate interface {
	AggregateID() string
	SetAggregateID(string)
	AggregateVersion() int
	DomainEvents() []Event
	On(event Event) error
}
