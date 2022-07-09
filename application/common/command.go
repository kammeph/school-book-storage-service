package common

type Command interface {
	AggregateID() string
}

type CommandModel struct {
	ID string `json:"aggregateId"`
}

func (c CommandModel) AggregateID() string {
	return c.ID
}
