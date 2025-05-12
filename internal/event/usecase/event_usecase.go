package usecase

import (
	"encoding/json"
	"fmt"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/event/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/dto/response"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"github.com/andrianprasetya/eventHub/internal/event/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	repositoryShared "github.com/andrianprasetya/eventHub/internal/shared/repository"
	responseDTO "github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	log "github.com/sirupsen/logrus"
)

type EventUsecase interface {
	Create(req request.CreateEventRequest, auth middleware.AuthUser, method string) error
	GetTags() ([]*response.EventTagListItemResponse, error)
	GetCategories() ([]*response.EventCategoryListItemResponse, error)
}

type eventUsecase struct {
	txManager         repositoryShared.TxManager
	eventRepo         repository.EventRepository
	eventTagRepo      repository.EventTagRepository
	eventCategoryRepo repository.EventCategoryRepository
	activityRepo      logRepository.LogActivityRepository
}

func NewEventUsecase(
	txManager repositoryShared.TxManager,
	eventRepo repository.EventRepository,
	eventTagRepo repository.EventTagRepository,
	eventCategoryRepo repository.EventCategoryRepository,
	activityRepo logRepository.LogActivityRepository,
) EventUsecase {
	return &eventUsecase{
		txManager:         txManager,
		eventRepo:         eventRepo,
		eventTagRepo:      eventTagRepo,
		eventCategoryRepo: eventCategoryRepo,
		activityRepo:      activityRepo}
}

func (u *eventUsecase) Create(req request.CreateEventRequest, auth middleware.AuthUser, method string) error {
	tx := u.txManager.Begin()
	event := &model.Event{
		ID:          utils.GenerateID(),
		TenantID:    auth.Tenant.ID,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      "published",
	}

	if errCreateEvent := u.eventRepo.Create(tx, event); errCreateEvent != nil {
		log.WithFields(log.Fields{
			"error": errCreateEvent,
		}).Error("failed to create event")
		return fmt.Errorf("something Went wrong")
	}

	userLog := responseDTO.EventLog{
		ID:          event.ID,
		TenantID:    event.TenantID,
		Description: event.Description,
		Location:    event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
	}

	userJSON, errMarshal := json.Marshal(userLog)
	if errMarshal != nil {
		return fmt.Errorf("error marshaling user")
	}

	helper.LogActivity(tx, u.activityRepo, auth.ID, method, string(userJSON), "event", event.ID)

	return nil
}

func (u *eventUsecase) GetTags() ([]*response.EventTagListItemResponse, error) {
	eventTags, err := u.eventTagRepo.GetAll()

	return mapper.FromEventTagToList(eventTags), err
}
func (u *eventUsecase) GetCategories() ([]*response.EventCategoryListItemResponse, error) {
	eventCategories, err := u.eventCategoryRepo.GetAll()

	return mapper.FromEventCategoryToList(eventCategories), err
}
