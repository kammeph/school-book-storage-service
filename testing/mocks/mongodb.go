package mocks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/infrastructure/mongodb"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockClient struct {
	mock.Mock
	database *MockDatabase
}

func NewMockClient() mongodb.Client {
	return &MockClient{database: NewMockDatabase()}
}

func (c *MockClient) Database(name string) mongodb.Database {
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

func (d *MockDatabase) Collection(name string) mongodb.Collection {
	return d.collection
}

type MockCollection struct {
	mock.Mock
}

func (c *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (mongodb.Cursor, error) {
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

func (c *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) mongodb.SingleResult {
	ret := c.Called(ctx, filter)

	var result *MockSingleResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(*MockSingleResult)
	} else {
		result = nil
	}
	return result
}

func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (mongodb.InsertResult, error) {
	ret := c.Called(ctx, document)

	var result mongodb.InsertResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(mongodb.InsertResult)
	}

	err := ret.Error(1)

	return result, err
}

func (c *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (mongodb.DeleteResult, error) {
	ret := c.Called(ctx, filter)

	var result mongodb.DeleteResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(mongodb.DeleteResult)
	}

	err := ret.Error(1)

	return result, err
}

func (c *MockCollection) UpdateOne(ctx context.Context, filter interface{}, document interface{}, opts ...*options.UpdateOptions) (mongodb.UpdateResult, error) {
	ret := c.Called(ctx, filter, document)

	var result mongodb.UpdateResult
	if ret.Get(0) != nil {
		result = ret.Get(0).(mongodb.UpdateResult)
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
