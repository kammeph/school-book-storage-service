package schoolapp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/domain/schooldomain"
)

type SchoolEventHandler struct {
	repository SchoolRepository
}

func NewSchoolEventHandler(repository SchoolRepository) application.EventHandler {
	return &SchoolEventHandler{repository}
}

func (h SchoolEventHandler) Handle(ctx context.Context, eventBytes []byte) {
	event := &domain.EventModel{}
	if err := json.Unmarshal(eventBytes, event); err != nil {
		log.Println(err.Error())
	}
	switch event.EventType() {
	case schooldomain.SchoolAdded:
		h.handleSchoolAdded(ctx, event)
	case schooldomain.SchoolDeactivated:
		h.handleSchoolDeactivated(ctx, event)
	case schooldomain.SchoolRenamed:
		h.handleSchoolRenamed(ctx, event)
	default:
		return
	}
}

func (h SchoolEventHandler) handleSchoolAdded(ctx context.Context, event domain.Event) {
	schoolAdded := schooldomain.SchoolAddedEvent{}
	if err := event.GetJsonData(&schoolAdded); err != nil {
		log.Println(err.Error())
	}
	school := schooldomain.NewSchoolProjection(schoolAdded.SchoolID, schoolAdded.Name)
	if err := h.repository.InsertSchool(ctx, school); err != nil {
		log.Println(err.Error())
	}
}

func (h SchoolEventHandler) handleSchoolDeactivated(ctx context.Context, event domain.Event) {
	eventData := schooldomain.SchoolDeactivatedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		log.Println(err.Error())
	}
	if err := h.repository.DeleteSchool(ctx, eventData.SchoolID); err != nil {
		log.Println(err.Error())
	}
}

func (h SchoolEventHandler) handleSchoolRenamed(ctx context.Context, event domain.Event) {
	eventData := schooldomain.SchoolRenamedEvent{}
	if err := event.GetJsonData(&eventData); err != nil {
		log.Println(err.Error())
	}
	if err := h.repository.UpdateSchoolName(ctx, eventData.SchoolID, eventData.Name); err != nil {
		log.Println(err.Error())
	}
}
