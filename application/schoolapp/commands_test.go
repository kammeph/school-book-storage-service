package schoolapp_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/application/schoolapp"
	"github.com/kammeph/school-book-storage-service/testing/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewSchoolCommandHandlers(t *testing.T) {
	store := &mocks.MockStore{}
	publisher := &mocks.MockEventPublisher{}
	commandHandlers := schoolapp.NewSchoolCommandHandlers(store, publisher)
	assert.NotNil(t, commandHandlers)
}

// func TestAddSchoolCommand(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		aggregateID   string
// 		schoolName    string
// 		loadedEvents  []domain.Event
// 		createdEvents []domain.Event
// 		loadErr       error
// 		saveErr       error
// 		publishErr    error
// 		expectedErr   error
// 	}{
// 		{
// 			name:          "add school",
// 			aggregateID:   "newSchool",
// 			schoolName:    "school",
// 			loadedEvents:  []domain.Event{},
// 			createdEvents: []domain.Event{},
// 			loadErr:       nil,
// 			saveErr:       nil,
// 			publishErr:    nil,
// 			expectedErr:   nil,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			store := mocks.NewMockStore()
// 			store.On("Load", context.Background(), test.aggregateID).Return(test.loadedEvents, test.loadErr)
// 			store.On("Save", context.Background(), test.createdEvents).Return(test.saveErr)
// 			publisher := mocks.NewMockEventPublisher()
// 			publisher.On("Publish", context.Background(), test.createdEvents).Return(test.publishErr)
// 			handler := schoolapp.NewAddStorageCommandHandler(store, publisher)
// 			command := schoolapp.AddSchoolCommand{
// 				CommandModel: application.CommandModel{ID: test.aggregateID},
// 				Name:         test.schoolName,
// 			}
// 			_, err := handler.Handle(context.Background(), command)
// 			assert.Equal(t, test.expectedErr, err)
// 		})
// 	}
// }
