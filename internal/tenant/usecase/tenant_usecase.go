package usecase

import (
	"context"
	eventRepository "github.com/andrianprasetya/eventHub/internal/event/repository"
	appErrors "github.com/andrianprasetya/eventHub/internal/shared/errors"
	sharedRepository "github.com/andrianprasetya/eventHub/internal/shared/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	userRepository "github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"runtime/debug"
	"time"
)

type TenantUsecase interface {
	RegisterTenant(request request.CreateTenantRequest) error
	UpdateInformation(id string, request request.UpdateInformationTenantRequest) error
	UpdateStatus(id string, request request.UpdateStatusTenantRequest) error
}

type tenantUsecase struct {
	txManager               sharedRepository.TxManager
	tenantRepo              repository.TenantRepository
	tenantSettingRepository repository.TenantSettingRepository
	subscriptionRepo        repository.SubscriptionRepository
	subscriptionPlanRepo    repository.SubscriptionPlanRepository
	userRepo                userRepository.UserRepository
	roleRepo                userRepository.RoleRepository
	eventTagRepo            eventRepository.EventTagRepository
	eventCategoryRepo       eventRepository.EventCategoryRepository
}

func NewTenantUsecase(
	txManager sharedRepository.TxManager,
	tenantRepo repository.TenantRepository,
	tenantSettingRepository repository.TenantSettingRepository,
	subscriptionRepo repository.SubscriptionRepository,
	subscriptionPlanRepo repository.SubscriptionPlanRepository,
	userRepo userRepository.UserRepository,
	roleRepo userRepository.RoleRepository,
	eventTagRepo eventRepository.EventTagRepository,
	eventCategoryRepo eventRepository.EventCategoryRepository,
) TenantUsecase {
	return &tenantUsecase{
		txManager:               txManager,
		tenantRepo:              tenantRepo,
		tenantSettingRepository: tenantSettingRepository,
		subscriptionRepo:        subscriptionRepo,
		subscriptionPlanRepo:    subscriptionPlanRepo,
		userRepo:                userRepo,
		roleRepo:                roleRepo,
		eventTagRepo:            eventTagRepo,
		eventCategoryRepo:       eventCategoryRepo,
	}
}

func (u *tenantUsecase) RegisterTenant(request request.CreateTenantRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var (
		subscriptionPlan *modelTenant.SubscriptionPlan
		role             *modelUser.Role
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		getPlan, err := u.subscriptionPlanRepo.GetByID(gctx, request.SubscriptionPlanID)
		if err != nil {
			return err
		}
		subscriptionPlan = getPlan
		return nil
	})

	g.Go(func() error {
		getRole, err := u.roleRepo.GetBySlug(gctx, "tenant-admin")
		if err != nil {
			return err
		}
		role = getRole
		return err
	})

	if err = g.Wait(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to fetch plan or role")
		return appErrors.WrapExpose(err, "failed to fetch plan or role", http.StatusInternalServerError)
	}

	tx := u.txManager.Begin(ctx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"errors": r,
				"stack":  string(debug.Stack()),
			}).Error("Failed to register tenant  (panic recovered)")
			err = appErrors.ErrInternalServer
		}
	}()

	hashedPassword, err := service.HashedPassword(request.Password)
	if err != nil {
		return appErrors.ErrInternalServer
	}

	tenant := service.MapTenantPayload(request)
	if err = u.tenantRepo.CreateWithTx(ctx, tx, tenant); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create tenant")
		return appErrors.ErrInternalServer
	}

	features, err := utils.ToStringJSON(subscriptionPlan.Feature)
	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to un-marshal feature")
		return appErrors.ErrInternalServer
	}

	tenantSettings := service.MapTenantSettingsPayload(features, tenant.ID)

	if err = u.tenantSettingRepository.CreateBulkWithTx(ctx, tx, tenantSettings); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to insert tenant settings")
		return appErrors.ErrInternalServer
	}

	eventCategories := service.BulkCategories(tenant.ID)
	eventTags := service.BulkTags(tenant.ID)

	if err = u.eventCategoryRepo.CreateBulkWithTx(ctx, tx, eventCategories); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to insert event tags")
		return appErrors.ErrInternalServer
	}

	if err = u.eventTagRepo.CreateBulkWithTx(ctx, tx, eventTags); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to insert event tags")
		return appErrors.ErrInternalServer
	}

	subscription := service.MapSubscriptionPayload(tenant.ID, subscriptionPlan.ID, subscriptionPlan.DurationDay)

	if err = u.subscriptionRepo.CreateWithTx(ctx, tx, subscription); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create subscription")
		return appErrors.ErrInternalServer
	}

	user := service.MapUserPayload(tenant.ID, role.ID, request.Name, request.Email, string(hashedPassword))

	if err = u.userRepo.CreateWithTx(ctx, tx, user); err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to create user")
		return appErrors.ErrInternalServer
	}

	if err = tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to commit transaction")
		tx.Rollback()
		return appErrors.ErrInternalServer
	}
	return err
}

func (u *tenantUsecase) UpdateInformation(id string, req request.UpdateInformationTenantRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tenant, err := u.tenantRepo.GetByID(ctx, id)

	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get tenant")
		return appErrors.ErrInternalServer
	}

	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.LogoUrl != nil {
		tenant.LogoUrl = *req.LogoUrl
	}

	if err := u.tenantRepo.Update(ctx, tenant); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to update tenant")
		return appErrors.ErrInternalServer
	}

	return err
}

func (u *tenantUsecase) UpdateStatus(id string, req request.UpdateStatusTenantRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tenant, err := u.tenantRepo.GetByID(ctx, id)

	if err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to get tenant")
		return appErrors.ErrInternalServer
	}

	tenant.IsActive = *req.IsActive

	if err := u.tenantRepo.Update(ctx, tenant); err != nil {
		log.WithFields(log.Fields{
			"errors": err,
		}).Error("failed to update tenant")
		return appErrors.ErrInternalServer
	}

	return err
}
