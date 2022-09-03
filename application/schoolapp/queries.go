package schoolapp

import (
	"context"

	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
)

type GetSchoolsQueryHandler struct {
	repository SchoolRepository
}

func NewGetSchoolsQueryHandler(repository SchoolRepository) GetSchoolsQueryHandler {
	return GetSchoolsQueryHandler{repository}
}

func (h GetSchoolsQueryHandler) Handle(ctx context.Context) ([]schooldomain.SchoolProjection, error) {
	return h.repository.GetSchools(ctx)
}

type GetSchoolByIDQuery struct {
	SchoolID string
}

func NewGetSchoolByIDQuery(id string) GetSchoolByIDQuery {
	return GetSchoolByIDQuery{SchoolID: id}
}

type GetSchoolByIDQueryHandler struct {
	repository SchoolRepository
}

func NewGetSchoolByIDQueryHandler(repository SchoolRepository) GetSchoolByIDQueryHandler {
	return GetSchoolByIDQueryHandler{repository}
}

func (h GetSchoolByIDQueryHandler) Handle(ctx context.Context, query GetSchoolByIDQuery) (schooldomain.SchoolProjection, error) {
	return h.repository.GetSchoolByID(ctx, query.SchoolID)
}
