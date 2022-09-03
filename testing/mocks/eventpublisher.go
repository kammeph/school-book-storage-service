package mocks

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

func (p *MockEventPublisher) Publish(ctx context.Context, events []domain.Event) error {
	ret := p.Called(ctx, events)
	err := ret.Error(0)
	return err
}
