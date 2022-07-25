package infrastructure

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type StorageWithBookRepository struct {
	collection Collection
}

func NewStorageWithBookRepository(db Client, dbName, tableName string) application.StorageWithBooksRepository {
	collection := db.Database(dbName).Collection(tableName)
	return &StorageWithBookRepository{collection}
}

func (c *StorageWithBookRepository) GetAllStoragesBySchoolID(ctx context.Context, schoolID string) ([]domain.StorageWithBooks, error) {
	filter := bson.D{{Key: "schoolId", Value: schoolID}}
	cursor, err := c.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	storages := []domain.StorageWithBooks{}
	if err := cursor.All(ctx, &storages); err != nil {
		return nil, err
	}
	return storages, nil
}

func (c *StorageWithBookRepository) GetStorageByID(ctx context.Context, schoolID, storageID string) (domain.StorageWithBooks, error) {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	result := c.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return domain.StorageWithBooks{}, result.Err()
	}

	storage := domain.StorageWithBooks{}
	if err := result.Decode(&storage); err != nil {
		return storage, err
	}
	return storage, nil
}

func (c *StorageWithBookRepository) GetStorageByName(ctx context.Context, schoolID, name string) (domain.StorageWithBooks, error) {
	filter := bson.D{{Key: "name", Value: name}}
	result := c.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return domain.StorageWithBooks{}, result.Err()
	}

	storage := domain.StorageWithBooks{}
	if err := result.Decode(&storage); err != nil {
		return storage, err
	}
	return storage, nil
}

func (c *StorageWithBookRepository) InsertStorage(ctx context.Context, storage domain.StorageWithBooks) error {
	doc := bson.D{
		{Key: "storageId", Value: storage.StorageID},
		{Key: "schoolId", Value: storage.SchoolID},
		{Key: "name", Value: storage.Name},
		{Key: "location", Value: storage.Location},
		{Key: "books", Value: storage.Books},
	}
	_, err := c.collection.InsertOne(ctx, doc)
	return err
}

func (c *StorageWithBookRepository) DeleteStorage(ctx context.Context, storageID string) error {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	_, err := c.collection.DeleteOne(ctx, filter)
	return err
}

func (c *StorageWithBookRepository) UpdateStorageName(ctx context.Context, storageID, name string) error {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: name}}}}
	_, err := c.collection.UpdateOne(ctx, filter, update)
	return err
}

func (c *StorageWithBookRepository) UpdateStorageLocation(ctx context.Context, storageID, location string) error {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "location", Value: location}}}}
	_, err := c.collection.UpdateOne(ctx, filter, update)
	return err
}
