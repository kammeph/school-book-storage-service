package mongodb

import (
	"context"

	"github.com/kammeph/school-book-storage-service/application/schoolapp"
	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
	"go.mongodb.org/mongo-driver/bson"
)

type SchoolRepository struct {
	collection Collection
}

func NewSchoolRepository(client Client, dbName, tableName string) schoolapp.SchoolRepository {
	collection := client.Database(dbName).Collection(tableName)
	return &SchoolRepository{collection}
}

func (r *SchoolRepository) GetSchools(ctx context.Context) ([]schooldomain.SchoolProjection, error) {
	filter := bson.D{}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	schools := []schooldomain.SchoolProjection{}
	if err := cursor.All(ctx, &schools); err != nil {
		return nil, err
	}
	return schools, nil
}

func (r *SchoolRepository) GetSchoolByID(ctx context.Context, schoolID string) (schooldomain.SchoolProjection, error) {
	filter := bson.D{{Key: "schoolId", Value: schoolID}}
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return schooldomain.SchoolProjection{}, result.Err()
	}

	school := schooldomain.SchoolProjection{}
	if err := result.Decode(&school); err != nil {
		return schooldomain.SchoolProjection{}, err
	}
	return school, nil
}

func (r *SchoolRepository) InsertSchool(ctx context.Context, school schooldomain.SchoolProjection) error {
	doc := bson.D{
		{Key: "schoolId", Value: school.SchoolID},
		{Key: "name", Value: school.Name},
	}
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *SchoolRepository) DeleteSchool(ctx context.Context, schoolID string) error {
	filter := bson.D{{Key: "schoolId", Value: schoolID}}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *SchoolRepository) UpdateSchoolName(ctx context.Context, schoolID, name string) error {
	filter := bson.D{{Key: "schoolId", Value: schoolID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: name}}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
