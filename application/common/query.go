package common

type QueryModel struct {
	ID string
}

type Query interface {
	AggregateID() string
}

func (q QueryModel) AggregateID() string {
	return q.ID
}
