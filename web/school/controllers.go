package school

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kammeph/school-book-storage-service/application/schoolapp"
	"github.com/kammeph/school-book-storage-service/web"
)

const (
	errAggregateID = "Aggregate ID is not specified"
	errSchoolID    = "School ID is not specified"
	errSchoolName  = "School name is not specified"
	errReason      = "Reason is not specified"
)

type SchoolController struct {
	commandHandlers schoolapp.SchoolCommandHandlers
	queryHandlers   schoolapp.SchoolQueryHandlers
}

func NewSchoolController(commandHandlers schoolapp.SchoolCommandHandlers, queryHandlers schoolapp.SchoolQueryHandlers) *SchoolController {
	return &SchoolController{commandHandlers, queryHandlers}
}

func (c SchoolController) AddSchool(w http.ResponseWriter, r *http.Request) {
	var command schoolapp.AddSchoolCommand
	json.NewDecoder(r.Body).Decode(&command)
	if command.AggregateID() == "" {
		web.HttpErrorResponse(w, errAggregateID)
		return
	}
	if command.Name == "" {
		web.HttpErrorResponse(w, errSchoolName)
		return
	}
	ctx := context.Background()
	defer ctx.Done()
	schoolID, err := c.commandHandlers.AddSchoolHandler.Handle(ctx, command)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	SchoolIDResponse(w, schoolID)
}

func (c SchoolController) DeactivateSchool(w http.ResponseWriter, r *http.Request) {
	var command schoolapp.DeactivateSchoolCommand
	json.NewDecoder(r.Body).Decode(&command)
	if command.AggregateID() == "" {
		web.HttpErrorResponse(w, errAggregateID)
		return
	}
	if command.SchoolID == "" {
		web.HttpErrorResponse(w, errSchoolID)
		return
	}
	if command.Reason == "" {
		web.HttpErrorResponse(w, errReason)
		return
	}
	ctx := context.Background()
	defer ctx.Done()
	if err := c.commandHandlers.DeactiveSchoolHandler.Handle(ctx, command); err != nil {
		web.HttpErrorResponse(w, err.Error())
	}
}

func (c SchoolController) RenameSchool(w http.ResponseWriter, r *http.Request) {
	var command schoolapp.RenameSchoolCommand
	json.NewDecoder(r.Body).Decode(&command)
	if command.AggregateID() == "" {
		web.HttpErrorResponse(w, errAggregateID)
		return
	}
	if command.SchoolID == "" {
		web.HttpErrorResponse(w, errSchoolID)
		return
	}
	if command.Name == "" {
		web.HttpErrorResponse(w, errSchoolName)
		return
	}
	if command.Reason == "" {
		web.HttpErrorResponse(w, errReason)
		return
	}
	ctx := context.Background()
	defer ctx.Done()
	if err := c.commandHandlers.RenameSchoolHandler.Handle(ctx, command); err != nil {
		web.HttpErrorResponse(w, err.Error())
	}
}

func (c SchoolController) GetSchools(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	schools, err := c.queryHandlers.GetSchoolsHandler.Handle(ctx)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	SchoolsResponse(w, schools)
}

func (c SchoolController) GetSchoolByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer ctx.Done()
	path := strings.Split(r.URL.Path, "/")
	schoolID := path[len(path)-1]
	query := schoolapp.NewGetSchoolByIDQuery(schoolID)
	school, err := c.queryHandlers.GetSchoolByIDHandler.Handle(ctx, query)
	if err != nil {
		web.HttpErrorResponse(w, err.Error())
		return
	}
	SchoolResponse(w, school)
}
