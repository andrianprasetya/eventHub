package usecase

import (
	"fmt"
	eventRepository "github.com/andrianprasetya/eventHub/internal/event/repository"
	repositoryShared "github.com/andrianprasetya/eventHub/internal/shared/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/service"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	userRepository "github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type TenantUsecase interface {
	RegisterTenant(request request.CreateTenantRequest) error
	Update(id string, request request.UpdateTenantRequest) error
}

type tenantUsecase struct {
	txManager               repositoryShared.TxManager
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
	txManager repositoryShared.TxManager,
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

func (u *tenantUsecase) RegisterTenant(request request.CreateTenantRequest) error {
	tx := u.txManager.Begin()

	var err error
	planCh := make(chan *modelTenant.SubscriptionPlanChannel)
	roleCh := make(chan *modelUser.RoleChannel)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"error": r,
			}).Error("Failed to create tenant  (panic recovered)")
			err = fmt.Errorf("something went wrong %w", r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tenant := &modelTenant.Tenant{
		ID:       utils.GenerateID(),
		Name:     request.Name,
		Email:    request.Email,
		LogoUrl:  request.LogoUrl,
		Domain:   utils.GenerateDomainName(request.Name),
		IsActive: 1,
	}

	go func() {
		plan, err := u.subscriptionPlanRepo.GetById(request.SubscriptionPlanID)
		planCh <- &modelTenant.SubscriptionPlanChannel{Plan: plan, Err: err}
	}()

	go func() {
		role, err := u.roleRepo.GetRole("tenant-admin")
		roleCh <- &modelUser.RoleChannel{Role: role, Err: err}
	}()

	resPlan := <-planCh
	if resPlan.Err != nil {
		log.WithFields(log.Fields{
			"error": resPlan.Err,
		}).Error("failed to get subscription plan")
		return fmt.Errorf("something Went wrong %w", resPlan.Err)
	}

	resRole := <-roleCh
	if resRole.Err != nil {
		log.WithFields(log.Fields{
			"error": resRole.Err,
		}).Error("failed to get role")
		return fmt.Errorf("something Went wrong")
	}

	if err = u.tenantRepo.CreateWithTx(tx, tenant); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create tenant")
		return fmt.Errorf("something Went wrong %w", err)
	}

	features, err := utils.ToStringJSON(resPlan.Plan.Feature)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to un-marshal feature")
		return fmt.Errorf("something Went wrong %w", err)
	}

	var tenantSettings []*modelTenant.TenantSetting
	for key, value := range features {
		strVal := fmt.Sprintf("%v", value)
		tenantSettings = append(tenantSettings, &modelTenant.TenantSetting{
			ID:       utils.GenerateID(),
			TenantID: tenant.ID,
			Key:      key,
			Value:    strVal,
		})
	}

	if err = u.tenantSettingRepository.CreateBulk(tx, tenantSettings); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to insert tenant settings")
		return fmt.Errorf("something Went wrong %w", err)
	}

	eventCategories := service.BulkCategories(tenant.ID)
	eventTags := service.BulkTags(tenant.ID)

	if err = u.eventCategoryRepo.CreateBulkWithTx(tx, eventCategories); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to insert event tags")
		return fmt.Errorf("something Went wrong %w", err)
	}

	if err = u.eventTagRepo.CreateBulkWithTx(tx, eventTags); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to insert event tags")
		return fmt.Errorf("something Went wrong %w", err)
	}

	var endDate *time.Time
	if resPlan.Plan.DurationDay != -1 {
		d := time.Now().AddDate(0, 0, resPlan.Plan.DurationDay)
		endDate = &d
	}

	subscription := &modelTenant.Subscription{
		ID:        utils.GenerateID(),
		TenantID:  tenant.ID,
		PlanID:    request.SubscriptionPlanID,
		StartDate: time.Now(),
		EndDate:   endDate,
		IsActive:  1,
	}

	if err = u.subscriptionRepo.Create(tx, subscription); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create subscription")
		return fmt.Errorf("something Went wrong %w", err)
	}

	user := &modelUser.User{
		ID:       utils.GenerateID(),
		TenantID: tenant.ID,
		RoleID:   resRole.Role.ID,
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		IsActive: 1,
	}
	if err = u.userRepo.Create(tx, user); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create user")
		return fmt.Errorf("something Went wrong %w", err)
	}

	err = tx.Commit().Error
	return err
}

func (u *tenantUsecase) Update(id string, req request.UpdateTenantRequest) error {
	tenant, err := u.tenantRepo.GetByID(id)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get tenant")
		return fmt.Errorf("something Went wrong %w", err)
	}

	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.LogoUrl != nil {
		tenant.LogoUrl = *req.LogoUrl
	}

	if err := u.tenantRepo.Update(tenant); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to update tenant")
		return fmt.Errorf("something Went wrong %w", err)
	}

	return nil
}
