package mocks

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s *MockStore) Load(ctx context.Context, aggregateID string) ([]domain.Event, error) {
	ret := s.Called(ctx, aggregateID)

	var events []domain.Event
	if ret.Get(0) != nil {
		events = ret.Get(0).([]domain.Event)
	} else {
		events = nil
	}

	err := ret.Error(1)

	return events, err
}

func (s *MockStore) Save(ctx context.Context, events []domain.Event) error {
	ret := s.Called(ctx, events)

	err := ret.Error(0)

	return err
}
