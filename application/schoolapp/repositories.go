package schoolapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
)

type SchoolRepository interface {
	GetSchools(ctx context.Context) ([]schooldomain.SchoolProjection, error)
	GetSchoolByID(ctx context.Context, schoolID string) (schooldomain.SchoolProjection, error)
	InsertSchool(ctx context.Context, school schooldomain.SchoolProjection) error
	DeleteSchool(ctx context.Context, schoolID string) error
	UpdateSchoolName(ctx context.Context, schoolID, name string) error
}
