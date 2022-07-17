package storage

import (
	"context"
	"fmt"

	application "github.com/kammeph/school-book-storage-service/application/storage"
	domain "github.com/kammeph/school-book-storage-service/domain/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageWithBookRepository struct {
	collection *mongo.Collection
}

func NewStorageWithBookRepository(db *mongo.Client, dbName, tableName string) application.StorageWithBooksRepository {
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
	if err := cursor.All(context.TODO(), &storages); err != nil {
		return nil, err
	}
	return storages, nil
}

func (c *StorageWithBookRepository) GetStorageByID(ctx context.Context, storageID string) (domain.StorageWithBooks, error) {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	result := c.collection.FindOne(ctx, filter)
	if result == nil {
		return domain.StorageWithBooks{}, fmt.Errorf("No storage with ID %s found", storageID)
	}

	storage := domain.StorageWithBooks{}
	if err := result.Decode(&storage); err != nil {
		return storage, err
	}
	return storage, nil
}

func (c *StorageWithBookRepository) GetStorageByName(ctx context.Context, name string) (domain.StorageWithBooks, error) {
	filter := bson.D{{Key: "name", Value: name}}
	result := c.collection.FindOne(ctx, filter)
	if result == nil {
		return domain.StorageWithBooks{}, fmt.Errorf("No storage with name %s found", name)
	}

	storage := domain.StorageWithBooks{}
	if err := result.Decode(&storage); err != nil {
		return storage, err
	}
	return storage, nil
}

func (c *StorageWithBookRepository) InsertStorage(ctx context.Context, storage domain.StorageWithBooks) {
	doc := bson.D{
		{Key: "storageId", Value: storage.StorageID},
		{Key: "schoolId", Value: storage.SchoolID},
		{Key: "name", Value: storage.Name},
		{Key: "location", Value: storage.Location},
		{Key: "books", Value: storage.Books},
	}
	_, err := c.collection.InsertOne(ctx, doc)
	if err != nil {
		fmt.Printf("Error while inserting storage with books: %s", err)
	}
}

func (c *StorageWithBookRepository) DeleteStorage(ctx context.Context, storageID string) {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	_, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("Error while deleting storage with ID %s", storageID)
	}
}

func (c *StorageWithBookRepository) UpdateStorageName(ctx context.Context, storageID, name string) {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: name}}}}
	_, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("Error while updating name from storage with ID %s", storageID)
	}
}

func (c *StorageWithBookRepository) UpdateStorageLocation(ctx context.Context, storageID, location string) {
	filter := bson.D{{Key: "storageId", Value: storageID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "location", Value: location}}}}
	_, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("Error while updating location from storage with ID %s", storageID)
	}
}
