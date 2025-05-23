package usecase

import (
	"context"
	"encoding/json"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	modelEvent "github.com/andrianprasetya/eventHub/internal/event/model"
	eventRepository "github.com/andrianprasetya/eventHub/internal/event/repository"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	sharedRepository "github.com/andrianprasetya/eventHub/internal/shared/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	tenantRepository "github.com/andrianprasetya/eventHub/internal/tenant/repository"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/response"
	"github.com/andrianprasetya/eventHub/internal/ticket/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

type TicketUsecase interface {
	Create(req request.CreateTicketRequest, auth middleware.AuthUser, url string) error
	GetAll(query request.TicketPaginateParams, auth middleware.AuthUser) ([]*response.TicketListItemResponse, int64, error)
	GetByID(id string) (*response.TicketResponse, error)
	Update(req request.UpdateTicketRequest, auth middleware.AuthUser, url string) error
	Delete(id string) error
}

type ticketUsecase struct {
	txManager         sharedRepository.TxManager
	ticketRepo        repository.TicketRepository
	eventRepo         eventRepository.EventRepository
	tenantSettingRepo tenantRepository.TenantSettingRepository
	activityRepo      logRepository.LogActivityRepository
}

func NewTicketUsecase(
	txManager sharedRepository.TxManager,
	ticketRepo repository.TicketRepository,
	eventRepo eventRepository.EventRepository,
	tenantSettingRepo tenantRepository.TenantSettingRepository,
	activityRepo logRepository.LogActivityRepository) TicketUsecase {
	return &ticketUsecase{
		txManager:         txManager,
		ticketRepo:        ticketRepo,
		eventRepo:         eventRepo,
		tenantSettingRepo: tenantSettingRepo,
		activityRepo:      activityRepo,
	}
}

func (u *ticketUsecase) Create(req request.CreateTicketRequest, auth middleware.AuthUser, url string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var (
		event           *modelEvent.Event
		unlimitedTicket string
		maxTicket       int
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		getEvent, err := u.eventRepo.GetByID(ctx, req.EventID)
		if err != nil {
			return err
		}
		event = getEvent
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
		getTenantSetting, err := u.tenantSettingRepo.GetByTenantID(gctx, auth.Tenant.ID, "max_tickets_per_event")
		if err != nil {
			return err
		}
		maxTicket, _ = strconv.Atoi(getTenantSetting.Value)
		return nil
	})

	if err = g.Wait(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to fetch event or tenant setting")
		return appErrors.WrapExpose(err, "failed to fetch plan or role", http.StatusInternalServerError)
	}

	tx := u.txManager.Begin(ctx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"error": r,
				"stack": string(debug.Stack()),
			}).Error("Recovered from panic in Create Ticket")
			err = appErrors.ErrInternalServer
		}
	}()
	eventTickets, err := service.MapEventTicket(event.ID, unlimitedTicket, req.Tickets, maxTicket)
	if err != nil {
		return appErrors.WrapExpose(err, "ticket quota Has been limit", http.StatusUnprocessableEntity)
	}

	if err = u.ticketRepo.CreateBulkWithTx(ctx, tx, eventTickets); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors":       err,
			"event_ticket": eventTickets,
		}).Error("failed to create ticket")
		return appErrors.ErrInternalServer
	}

	if event.IsTicket == 0 {
		event.IsTicket = 1
		if err := u.eventRepo.UpdateWithTx(ctx, tx, event); err != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"errors": err,
				"event":  event,
			}).Error("failed to update event")
			return appErrors.ErrInternalServer
		}
	}
	eventJSON, errMarshal := json.Marshal(eventTickets)
	if errMarshal != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors":    err,
			"event_log": event,
		}).Error("failed to create event session")
		return appErrors.ErrInternalServer
	}

	err = tx.Commit().Error
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to commit transaction")
		tx.Rollback()
		return appErrors.ErrInternalServer
	}

	helper.LogActivity(u.activityRepo, ctx, auth.Tenant.ID, auth.ID, url, "Create Event Ticket", string(eventJSON), "event", event.ID)
	return nil
}

func (u *ticketUsecase) GetAll(query request.TicketPaginateParams, auth middleware.AuthUser) ([]*response.TicketListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	getTickets, total, err := u.ticketRepo.GetAll(ctx, query, &auth.Tenant.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"query_params": query,
			"tenant_id":    auth.Tenant.ID,
		}).Error("failed to create event ticket")
		return nil, 0, appErrors.ErrInternalServer
	}

	return mapper.FromTicketToList(getTickets), total, nil
}

func (u *ticketUsecase) GetByID(id string) (*response.TicketResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	getTicket, err := u.ticketRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to create event ticket")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromTicketModel(getTicket), nil
}

func (u *ticketUsecase) Update(req request.UpdateTicketRequest, auth middleware.AuthUser, url string) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	getTicket, err := u.ticketRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to create event ticket")
		return appErrors.ErrInternalServer
	}*/

	return nil
}

func (u *ticketUsecase) Delete(id string) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()*/
	return nil
}
