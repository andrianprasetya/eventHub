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
	Create(req request.CreateEventRequest, auth middleware.AuthUser, url string) (*response.EventResponse, error)
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

func (u *eventUsecase) Create(req request.CreateEventRequest, auth middleware.AuthUser, url string) (*response.EventResponse, error) {
	tx := u.txManager.Begin()

	var err error

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"error": r,
			}).Error("Failed to create event  (panic recovered)")
			err = fmt.Errorf("something went wrong")
		} else if err != nil {
			tx.Rollback()
		}
	}()

	event := &model.Event{
		ID:          utils.GenerateID(),
		Title:       req.Title,
		TenantID:    auth.Tenant.ID,
		CategoryID:  req.CategoryID,
		Tags:        req.Tags,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      "draft",
	}

	if err = u.eventRepo.Create(tx, event); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create event")
		return &response.EventResponse{}, err
	}

	if err = u.eventCategoryRepo.AddCategoryToEventWithTx(tx, event.ID, event); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to add category to event")
		return &response.EventResponse{}, err
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
		return &response.EventResponse{}, err
	}
	err = tx.Commit().Error
	if err == nil {
		helper.LogActivity(u.activityRepo, auth.ID, url, "Create Event", string(userJSON), "event", event.ID)
	}

	return mapper.FromUserModel(event), err
}

func (u *eventUsecase) GetTags() ([]*response.EventTagListItemResponse, error) {
	eventTags, err := u.eventTagRepo.GetAll()

	return mapper.FromEventTagToList(eventTags), err
}
func (u *eventUsecase) GetCategories() ([]*response.EventCategoryListItemResponse, error) {
	eventCategories, err := u.eventCategoryRepo.GetAll()

	return mapper.FromEventCategoryToList(eventCategories), err
}
