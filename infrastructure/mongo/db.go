package mongo

import (
	"context"
	"fmt"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	user     = utils.GetenvOrFallback("PG_USER", "mongo")
	password = utils.GetenvOrFallback("PG_PASSWORD", "mongo")
	host     = utils.GetenvOrFallback("PG_HOST", "localhost")
	port     = utils.GetenvOrFallback("PG_PORT", "27017")
)

func NewDB() *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged to mongo db.")
	return client
}
