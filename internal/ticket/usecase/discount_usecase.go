package usecase

import (
	"context"
	"encoding/json"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	eventRepository "github.com/andrianprasetya/eventHub/internal/event/repository"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	"github.com/andrianprasetya/eventHub/internal/shared/helper"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/mapper"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/response"
	"github.com/andrianprasetya/eventHub/internal/ticket/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

type DiscountUsecase interface {
	CreateBulk(req request.CreateDiscountRequest, auth middleware.AuthUser, url string) error
	GetAll(query request.DiscountPaginateParams, userAuth middleware.AuthUser) ([]*response.DiscountListItemResponse, int64, error)
	GetByID(id string) (*response.DiscountResponse, error)
}

type discountUsecase struct {
	discountRepo repository.DiscountRepository
	eventRepo    eventRepository.EventRepository
	activityRepo logRepository.LogActivityRepository
}

func NewDiscountUsecase(
	discountRepo repository.DiscountRepository,
	eventRepo eventRepository.EventRepository,
	activityRepo logRepository.LogActivityRepository) DiscountUsecase {
	return &discountUsecase{
		discountRepo: discountRepo,
		eventRepo:    eventRepo,
		activityRepo: activityRepo,
	}
}

func (u *discountUsecase) CreateBulk(req request.CreateDiscountRequest, auth middleware.AuthUser, url string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	getEvent, err := u.eventRepo.GetByID(ctx, req.EventID)

	if err != nil {
		log.WithFields(log.Fields{
			"errors":   err,
			"event_id": req.EventID,
		}).Error("failed to get role")
		return appErrors.ErrInternalServer
	}

	discounts := service.MapDiscountsPayload(getEvent.ID, req.Discounts)

	discountJSON, errMarshal := json.Marshal(discounts)
	if errMarshal != nil {
		log.WithFields(log.Fields{
			"errors":    err,
			"event_log": discounts,
		}).Error("failed to create event session")
		return appErrors.ErrInternalServer
	}

	if err = u.discountRepo.CreateBulk(ctx, discounts); err != nil {
		log.WithFields(log.Fields{
			"errors":    err,
			"discounts": discounts,
		}).Error("failed to create discount ticket")
		return appErrors.ErrInternalServer
	}
	helper.LogActivity(u.activityRepo, ctx, auth.Tenant.ID, auth.ID, url, "Create Event Discount", string(discountJSON), "event_discount", getEvent.ID)
	return err
}

func (u *discountUsecase) GetAll(query request.DiscountPaginateParams, auth middleware.AuthUser) ([]*response.DiscountListItemResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	getDiscounts, total, err := u.discountRepo.GetAll(ctx, query, &auth.Tenant.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"errors":       err,
			"query_params": query,
			"tenant_id":    auth.Tenant.ID,
		}).Error("failed to create discount ticket")
		return nil, 0, appErrors.ErrInternalServer
	}
	return mapper.FromDiscountToList(getDiscounts), total, nil
}

func (u *discountUsecase) GetByID(id string) (*response.DiscountResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	getDiscount, err := u.discountRepo.GetByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
			"id":     id,
		}).Error("failed to create discount ticket")
		return nil, appErrors.ErrInternalServer
	}
	return mapper.FromDiscountModel(getDiscount), nil
}
