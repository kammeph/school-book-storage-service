package common

type Aggregate interface {
	On(event Event)
}
