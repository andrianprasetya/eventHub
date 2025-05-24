package usecase

import (
	"context"
	"encoding/json"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/event/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/dto/response"
	"github.com/andrianprasetya/eventHub/internal/event/repository"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	sharedRepository "github.com/andrianprasetya/eventHub/internal/shared/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	tenantRepository "github.com/andrianprasetya/eventHub/internal/tenant/repository"
	ticketRepository "github.com/andrianprasetya/eventHub/internal/ticket/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

type EventUsecase interface {
	Create(req request.CreateEventRequest, auth middleware.AuthUser, url string) (*response.EventResponse, error)
	GetTags(query request.EventTagPaginateRequest, tenantID *string) ([]*response.EventTagListItemResponse, int64, error)
	GetCategories(query request.EventCategoryPaginateRequest, tenantID *string) ([]*response.EventCategoryListItemResponse, int64, error)
	GetAll(query request.EventPaginateRequest, tenantID *string) ([]*response.EventListItemResponse, int64, error)
	GetByID(id string) (*response.EventResponse, error)
	Update(req request.UpdateEventRequest, id string) error
	Delete(id string) error
}

type eventUsecase struct {
	txManager         sharedRepository.TxManager
	tenantSettingRepo tenantRepository.TenantSettingRepository
	eventRepo         repository.EventRepository
	eventTagRepo      repository.EventTagRepository
	eventCategoryRepo repository.EventCategoryRepository
	eventSessionRepo  repository.EventSessionRepository
	ticketRepo        ticketRepository.TicketRepository
	discountRepo      ticketRepository.DiscountRepository
	activityRepo      logRepository.LogActivityRepository
}

func NewEventUsecase(
	txManager sharedRepository.TxManager,
	tenantSettingRepo tenantRepository.TenantSettingRepository,
	eventRepo repository.EventRepository,
	eventTagRepo repository.EventTagRepository,
	eventCategoryRepo repository.EventCategoryRepository,
	eventSessionRepo repository.EventSessionRepository,
	ticketRepo ticketRepository.TicketRepository,
	discountRepo ticketRepository.DiscountRepository,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	var (
		maxEvent          int
		countEventCreated int
		unlimitedEvent    string
		unlimitedTicket   string
		maxTicket         int
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		getTenantSetting, err := u.tenantSettingRepo.GetByTenantID(gctx, auth.Tenant.ID, "unlimited_events")
		if err != nil {
			return err
		}
		unlimitedEvent = getTenantSetting.Value
		return nil
	})

	g.Go(func() error {
		getTenantSetting, err := u.tenantSettingRepo.GetByTenantID(gctx, auth.Tenant.ID, "unlimited_tickets_per_event")
		if err != nil {
			return err
		}
		unlimitedTicket = getTenantSetting.Value
		return nil
	})

	g.Go(func() error {
		getTenantSetting, err := u.tenantSettingRepo.GetByTenantID(gctx, auth.Tenant.ID, "max_events")
		if err != nil {
			return err
		}
		maxEvent, _ = strconv.Atoi(getTenantSetting.Value)
		return nil
	})
	g.Go(func() error {
		getTenantSetting, err := u.tenantSettingRepo.GetByTenantID(gctx, auth.Tenant.ID, "max_tickets_per_event")
		if err != nil {
			return err
		}
		maxTicket, _ = strconv.Atoi(getTenantSetting.Value)
		return nil
	})

	g.Go(func() error {
		getEventHasCreated, err := u.eventRepo.CountCreatedEvent(gctx, auth.Tenant.ID)
		if err != nil {
			return err
		}
		countEventCreated = getEventHasCreated
		return nil
	})
	if err = g.Wait(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to fetch tenant")
		return nil, appErrors.WrapExpose(err, "failed to fetch plan or role", http.StatusInternalServerError)
	}

	tx := u.txManager.Begin(ctx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"error": r,
				"stack": string(debug.Stack()),
			}).Error("Recovered from panic in Create Event")
			err = appErrors.ErrInternalServer
		}
	}()

	if unlimitedEvent == "false" {
		if err = service.CheckMaxEventCanCreated(countEventCreated, maxEvent); err != nil {
			return nil, appErrors.WrapExpose(err, "Created event quota Has been limit", http.StatusUnprocessableEntity)
		}
	}
	event, err := service.MapEventPayload(auth, req)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":  err,
			"auth":    auth,
			"request": req,
		})
		return nil, appErrors.ErrInternalServer
	}

	if err = u.eventRepo.Create(ctx, tx, event); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
			"event":  event,
		}).Error("failed to create event")
		return nil, appErrors.ErrInternalServer
	}

	if err = u.eventCategoryRepo.AddCategoryToEventWithTx(ctx, tx, event); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
			"event":  event,
		}).Error("failed to add category to event")
		return nil, appErrors.ErrInternalServer
	}

	eventTickets, err := service.MapEventTicket(event.ID, unlimitedTicket, req.Tickets, maxTicket)
	if err != nil {
		tx.Rollback()
		return nil, appErrors.WrapExpose(err, "ticket quota Has been limit", http.StatusUnprocessableEntity)
	}

	if err = u.ticketRepo.CreateBulkWithTx(ctx, tx, eventTickets); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors":       err,
			"event_ticket": eventTickets,
		}).Error("failed to create ticket")
		return nil, appErrors.ErrInternalServer
	}

	discounts := service.MapDiscountsPayload(event.ID, req.Discounts)
	if err = u.discountRepo.CreateBulkWithTx(ctx, tx, discounts); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors":    err,
			"discounts": discounts,
		}).Error("failed to create discount ticket")
		return nil, appErrors.ErrInternalServer
	}

	sessions := service.MapEventServicesPayload(event.ID, req.Sessions)

	if err = u.eventSessionRepo.CreateBulkWithTx(ctx, tx, sessions); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors":   err,
			"sessions": sessions,
		}).Error("failed to create event session")
		return nil, appErrors.ErrInternalServer
	}

	userJSON, errMarshal := json.Marshal(event)
	if errMarshal != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors":    err,
			"event_log": event,
		}).Error("failed to create event session")
		return nil, appErrors.ErrInternalServer
	}
	if err = tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to commit transaction")
		tx.Rollback()
		return nil, appErrors.ErrInternalServer
	}
	helper.LogActivity(u.activityRepo, ctx, auth.Tenant.ID, auth.ID, url, "Create Event", string(userJSON), "event", event.ID)

	return mapper.FromEventModel(event), nil
}

func (u *eventUsecase) GetTags(query request.EventTagPaginateRequest, tenantID *string) ([]*response.EventTagListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventTags, total, err := u.eventTagRepo.GetAll(ctx, query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"query_params": query,
			"tenant_id":    tenantID,
		}).Error("failed to get tags")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromEventTagToList(eventTags), total, err
}
func (u *eventUsecase) GetCategories(query request.EventCategoryPaginateRequest, tenantID *string) ([]*response.EventCategoryListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventCategories, total, err := u.eventCategoryRepo.GetAll(ctx, query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"query_params": query,
			"tenant_id":    tenantID,
		}).Error("failed to get categories")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromEventCategoryToList(eventCategories), total, err
}

func (u *eventUsecase) GetAll(query request.EventPaginateRequest, tenantID *string) ([]*response.EventListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	events, total, err := u.eventRepo.GetAll(ctx, query, tenantID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"query_params": query,
			"tenant_id":    tenantID,
		}).Error("failed to get categories")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromEventToList(events), total, err
}

func (u *eventUsecase) GetByID(id string) (*response.EventResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	event, err := u.eventRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to get categories")
		return nil, appErrors.ErrInternalServer
	}

	return mapper.FromEventModel(event), err
}

func (u *eventUsecase) Update(req request.UpdateEventRequest, id string) error {
	return nil
}

func (u *eventUsecase) Delete(id string) error {
	return nil
}
