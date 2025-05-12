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

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.WithFields(log.Fields{
				"error": r,
			}).Error("Failed to create tenant  (panic recovered)")
			err = fmt.Errorf("something went wrong")
		} else if err != nil {
			tx.Rollback()
		}
	}()

	subscriptionPlan, err := u.subscriptionPlanRepo.GetById(request.SubscriptionPlanID)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get subscription plan")
		return fmt.Errorf("something Went wrong")
	}

	role, err := u.roleRepo.GetRole("tenant-admin")

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to get role")
		return fmt.Errorf("something Went wrong")
	}

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

	if err = u.tenantRepo.Create(tx, tenant); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create tenant")
		return fmt.Errorf("something Went wrong")
	}

	features, err := utils.ToStringJSON(subscriptionPlan.Feature)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to un-marshal feature")
		return fmt.Errorf("something Went wrong")
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

	if err := u.tenantSettingRepository.CreateBulk(tx, tenantSettings); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to insert tenant settings")
		return fmt.Errorf("something went wrong")
	}

	eventCategories := service.BulkCategories(tenant.ID)
	eventTags := service.BulkTags(tenant.ID)

	if err := u.eventCategoryRepo.CreateBulkWithTx(tx, eventCategories); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to insert event tags")
		return fmt.Errorf("something went wrong")
	}

	if err := u.eventTagRepo.CreateBulkWithTx(tx, eventTags); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to insert event tags")
		return fmt.Errorf("something went wrong")
	}

	var endDate *time.Time
	if subscriptionPlan.DurationDay != -1 {
		d := time.Now().AddDate(0, 0, subscriptionPlan.DurationDay)
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
		return fmt.Errorf("something Went wrong")
	}

	user := &modelUser.User{
		ID:       utils.GenerateID(),
		TenantID: tenant.ID,
		RoleID:   role.ID,
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		IsActive: 1,
	}
	if err := u.userRepo.Create(tx, user); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to create user")
		return fmt.Errorf("something Went wrong")
	}

	return tx.Commit().Error
}

func (u *tenantUsecase) Update(id string, req request.UpdateTenantRequest) error {
	tenant, errGetTenant := u.tenantRepo.GetByID(id)

	if errGetTenant != nil {
		log.WithFields(log.Fields{
			"error": errGetTenant,
		}).Error("failed to get tenant")
		return fmt.Errorf("something Went wrong")
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
		return fmt.Errorf("something Went wrong")
	}

	return nil
}
