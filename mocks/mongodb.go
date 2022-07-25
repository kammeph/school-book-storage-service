package mocks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/infrastructure"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockClient struct {
	mock.Mock
	database *MockDatabase
}

func NewMockClient() infrastructure.Client {
	return &MockClient{database: NewMockDatabase()}
}

func (c *MockClient) Database(name string) infrastructure.Database {
	return c.database
}

func (c *MockClient) Disconnect(ctx context.Context) error {
	return nil
}

type MockDatabase struct {
	mock.Mock
	collection *MockCollection
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{collection: &MockCollection{}}
}

func (d *MockDatabase) Collection(name string) infrastructure.Collection {
	return d.collection
}

type MockCollection struct {
	mock.Mock
}

func (c *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (infrastructure.Cursor, error) {
	ret := c.Called(ctx, filter)

	var data []byte
	if ret.Get(0) != nil {
		storages := ret.Get(0).([]domain.StorageWithBooks)
		data, _ = json.Marshal(storages)
	} else {
		data = nil
	}

	err := ret.Error(1)

	return &MockCursor{data: data}, err
}

func (c *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) infrastructure.SingleResult {
	ret := c.Called(ctx, filter)

	var result *MockSingleResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(*MockSingleResult)
	} else {
		result = nil
	}
	return result
}

func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (infrastructure.InsertResult, error) {
	ret := c.Called(ctx, document)

	var result infrastructure.InsertResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(infrastructure.InsertResult)
	}

	err := ret.Error(1)

	return result, err
}

func (c *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (infrastructure.DeleteResult, error) {
	ret := c.Called(ctx, filter)

	var result infrastructure.DeleteResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(infrastructure.DeleteResult)
	}

	err := ret.Error(1)

	return result, err
}

func (c *MockCollection) UpdateOne(ctx context.Context, filter interface{}, document interface{}, opts ...*options.UpdateOptions) (infrastructure.UpdateResult, error) {
	ret := c.Called(ctx, filter, document)

	var result infrastructure.UpdateResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(infrastructure.UpdateResult)
	}

	err := ret.Error(1)

	return result, err
}

type MockCursor struct {
	data []byte
}

func (c *MockCursor) All(ctx context.Context, result interface{}) error {
	if c.data == nil {
		return errors.New("mock-cursor-error")
	}
	json.Unmarshal(c.data, result)
	return nil
}

type MockSingleResult struct {
	Storage *domain.StorageWithBooks
	Error   error
}

func (r *MockSingleResult) Decode(result interface{}) error {
	if r.Storage == nil {
		return errors.New("mock-decode-error")
	}
	data, _ := json.Marshal(r.Storage)
	json.Unmarshal(data, result)
	return nil
}

func (r *MockSingleResult) Err() error {
	return r.Error
}
