package mongodb

import (
	"context"
	"fmt"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	user     = utils.GetenvOrFallback("MONGO_USER", "mongo")
	password = utils.GetenvOrFallback("MONGO_PASSWORD", "mongo")
	host     = utils.GetenvOrFallback("MONGO_HOST", "localhost")
	port     = utils.GetenvOrFallback("MONGO_PORT", "27017")
)

type ClientWrapper struct {
	*mongo.Client
}

func (c *ClientWrapper) Database(name string) Database {
	db := c.Client.Database(name)
	return &DatabaseWrapper{Database: db}
}

func (c *ClientWrapper) Disconnect(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}

type DatabaseWrapper struct {
	*mongo.Database
}

func (d *DatabaseWrapper) Collection(name string) Collection {
	collection := d.Database.Collection(name)
	return &CollectionWrapper{Collection: collection}
}

type CollectionWrapper struct {
	*mongo.Collection
}

func (c *CollectionWrapper) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	return c.Collection.Find(ctx, filter, opts...)
}

func (c *CollectionWrapper) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	return c.Collection.FindOne(ctx, filter, opts...)
}

func (c *CollectionWrapper) InsertOne(ctx context.Context, filter interface{}, opts ...*options.InsertOneOptions) (InsertResult, error) {
	return c.Collection.InsertOne(ctx, filter, opts...)
}

func (c *CollectionWrapper) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	return c.Collection.DeleteOne(ctx, filter, opts...)
}

func (c *CollectionWrapper) UpdateOne(ctx context.Context, filter interface{}, document interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	return c.Collection.UpdateOne(ctx, filter, document, opts...)
}

func NewMongoClient() Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged to mongo db.")
	return &ClientWrapper{client}
}

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
