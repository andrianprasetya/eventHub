package usecase

import (
	"encoding/json"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/event/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/dto/response"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"github.com/andrianprasetya/eventHub/internal/event/repository"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	repositoryShared "github.com/andrianprasetya/eventHub/internal/shared/repository"
	responseDTO "github.com/andrianprasetya/eventHub/internal/shared/response"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	repositoruTenant "github.com/andrianprasetya/eventHub/internal/tenant/repository"
	modelTicket "github.com/andrianprasetya/eventHub/internal/ticket/model"
	repositoryTicket "github.com/andrianprasetya/eventHub/internal/ticket/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type EventUsecase interface {
	Create(req request.CreateEventRequest, auth middleware.AuthUser, url string) (*response.EventResponse, error)
	GetTags(query request.EventTagPaginateRequest, tenantID *string) ([]*response.EventTagListItemResponse, int64, error)
	GetCategories(query request.EventCategoryPaginateRequest, tenantID *string) ([]*response.EventCategoryListItemResponse, int64, error)
	GetAll(query request.EventPaginateRequest, tenantID *string) ([]*response.EventListItemResponse, int64, error)
	GetByID(id string) (*response.EventResponse, error)
}

type eventUsecase struct {
	txManager         repositoryShared.TxManager
	tenantSettingRepo repositoruTenant.TenantSettingRepository
	eventRepo         repository.EventRepository
	eventTagRepo      repository.EventTagRepository
	eventCategoryRepo repository.EventCategoryRepository
	eventSessionRepo  repository.EventSessionRepository
	ticketRepo        repositoryTicket.TicketRepository
	discountRepo      repositoryTicket.DiscountRepository
	activityRepo      logRepository.LogActivityRepository
}

func NewEventUsecase(
	txManager repositoryShared.TxManager,
	tenantSettingRepo repositoruTenant.TenantSettingRepository,
	eventRepo repository.EventRepository,
	eventTagRepo repository.EventTagRepository,
	eventCategoryRepo repository.EventCategoryRepository,
	eventSessionRepo repository.EventSessionRepository,
	ticketRepo repositoryTicket.TicketRepository,
	discountRepo repositoryTicket.DiscountRepository,
	activityRepo logRepository.LogActivityRepository,
) EventUsecase {
	return &eventUsecase{
		txManager:         txManager,
		tenantSettingRepo: tenantSettingRepo,
		eventRepo:         eventRepo,
		eventTagRepo:      eventTagRepo,
		eventCategoryRepo: eventCategoryRepo,
		eventSessionRepo:  eventSessionRepo,
		ticketRepo:        ticketRepo,
		discountRepo:      discountRepo,
		activityRepo:      activityRepo}
}

func (u *eventUsecase) Create(req request.CreateEventRequest, auth middleware.AuthUser, url string) (*response.EventResponse, error) {

	tx := u.txManager.Begin()

	var err error

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"errors": r,
			}).Error("Failed to create event  (panic recovered)")
			err = appErrors.ErrInternalServer
		} else if err != nil {
			tx.Rollback()
		}
	}()
	tenantSettings, err := u.tenantSettingRepo.GetByTenantID(auth.Tenant.ID, "max_events")
	countEventCreated := u.eventRepo.CountCreatedEvent(auth.Tenant.ID)
	if err != nil {
		return nil, appErrors.ErrInternalServer
	}

	if err = service.CheckMaxEventCanCreated(countEventCreated, tenantSettings); err != nil {
		return nil, appErrors.WrapExpose(err, "Created event quota Has been limit", http.StatusUnprocessableEntity)
	}

	event := &model.Event{
		ID:          utils.GenerateID(),
		Title:       req.Title,
		TenantID:    auth.Tenant.ID,
		CategoryID:  req.CategoryID,
		EventType:   req.EventType,
		Tags:        req.Tags,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedBy:   auth.ID,
		IsTicket:    req.IsTicket,
		Status:      req.Status,
	}

	if err = u.eventRepo.Create(tx, event); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create event")
		return nil, appErrors.ErrInternalServer
	}

	if err = u.eventCategoryRepo.AddCategoryToEventWithTx(tx, event.ID, event); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to add category to event")
		return nil, appErrors.ErrInternalServer
	}

	var eventTickets []*modelTicket.EventTicket
	for _, ticket := range req.Tickets {
		eventTickets = append(eventTickets, &modelTicket.EventTicket{
			ID:         utils.GenerateID(),
			EventID:    event.ID,
			TicketType: ticket.Type,
			Price:      ticket.Price,
			Quantity:   ticket.Quantity,
		})
	}

	if err = u.ticketRepo.CreateBulkWithTx(tx, eventTickets); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create ticket")
		return nil, appErrors.ErrInternalServer
	}

	var discounts []*modelTicket.Discount
	for _, discount := range req.Discounts {
		discounts = append(discounts, &modelTicket.Discount{
			ID:                 utils.GenerateID(),
			EventID:            event.ID,
			Code:               discount.Code,
			DiscountPercentage: discount.DiscountPercentage,
			StartDate:          discount.StartDate,
			EndDate:            discount.EndDate,
		})
	}

	if err = u.discountRepo.CreateBulkWithTx(tx, discounts); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create discount ticket")
		return nil, appErrors.ErrInternalServer
	}

	var sessions []*model.EventSession
	for _, session := range req.Sessions {
		sessions = append(sessions, &model.EventSession{
			ID:            utils.GenerateID(),
			EventID:       event.ID,
			Title:         session.Title,
			StartDateTime: session.StartDateTime,
			EndDateTime:   session.EndDateTime,
		})
	}

	if err = u.eventSessionRepo.CreateBulkWithTx(tx, sessions); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create event session")
		return nil, appErrors.ErrInternalServer
	}

	userLog := responseDTO.EventLog{
		ID:          event.ID,
		TenantID:    event.TenantID,
		Description: event.Description,
		Location:    *event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
	}

	userJSON, errMarshal := json.Marshal(userLog)
	if errMarshal != nil {
		return nil, appErrors.ErrInternalServer
	}
	err = tx.Commit().Error
	if err == nil {
		helper.LogActivity(u.activityRepo, auth.Tenant.ID, auth.ID, url, "Create Event", string(userJSON), "event", event.ID)
	}

	return mapper.FromEventModel(event), err
}

func (u *eventUsecase) GetTags(query request.EventTagPaginateRequest, tenantID *string) ([]*response.EventTagListItemResponse, int64, error) {
	eventTags, total, err := u.eventTagRepo.GetAll(query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get tags")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromEventTagToList(eventTags), total, err
}
func (u *eventUsecase) GetCategories(query request.EventCategoryPaginateRequest, tenantID *string) ([]*response.EventCategoryListItemResponse, int64, error) {
	eventCategories, total, err := u.eventCategoryRepo.GetAll(query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get categories")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromEventCategoryToList(eventCategories), total, err
}

func (u *eventUsecase) GetAll(query request.EventPaginateRequest, tenantID *string) ([]*response.EventListItemResponse, int64, error) {
	events, total, err := u.eventRepo.GetAll(query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get categories")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromEventToList(events), total, err
}

func (u *eventUsecase) GetByID(id string) (*response.EventResponse, error) {
	event, err := u.eventRepo.GetByID(id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get categories")
		return nil, appErrors.ErrInternalServer
	}

	return mapper.FromEventModel(event), err
}
