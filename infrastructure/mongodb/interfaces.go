package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Database(name string) Database
	Disconnect(ctx context.Context) error
}

type Database interface {
	Collection(name string) Collection
}

type Collection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (InsertResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error)
}

type Cursor interface {
	All(ctx context.Context, result interface{}) error
}

type SingleResult interface {
	Decode(v interface{}) error
	Err() error
}

type InsertResult interface {
}

type DeleteResult interface {
}

type UpdateResult interface {
}
