package usecase

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"github.com/andrianprasetya/eventHub/internal/tenant/repository"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	userRepository "github.com/andrianprasetya/eventHub/internal/user/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type TenantUsecase interface {
	RegisterTenant(request request.CreateTenantRequest) error
}

type tenantUsecase struct {
	tenantRepo           repository.TenantRepository
	subscriptionRepo     repository.SubscriptionRepository
	subscriptionPlanRepo repository.SubscriptionPlanRepository
	userRepo             userRepository.UserRepository
	roleRepo             userRepository.RoleRepository
}

func NewTenantUsecase(
	tenantRepo repository.TenantRepository,
	subscriptionRepo repository.SubscriptionRepository,
	subscriptionPlanRepo repository.SubscriptionPlanRepository,
	userRepo userRepository.UserRepository,
	roleRepo userRepository.RoleRepository) TenantUsecase {
	return &tenantUsecase{
		tenantRepo:           tenantRepo,
		subscriptionRepo:     subscriptionRepo,
		subscriptionPlanRepo: subscriptionPlanRepo,
		userRepo:             userRepo,
		roleRepo:             roleRepo,
	}
}

func (u tenantUsecase) RegisterTenant(request request.CreateTenantRequest) error {
	subscriptionPlan, errGetPlan := u.subscriptionPlanRepo.Get(request.SubscriptionPlanID)

	role, errGetRole := u.roleRepo.GetRole("admin")

	if errGetPlan != nil {
		log.WithFields(log.Fields{
			"error": errGetPlan,
		}).Error("failed to get subscription plan")
		return fmt.Errorf("something Went wrong")
	}
	if errGetRole != nil {
		log.WithFields(log.Fields{
			"error": errGetRole,
		}).Error("failed to get role")
		return fmt.Errorf("something Went wrong")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("TENANT_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	tenant := &modelTenant.Tenant{
		ID:       utils.GenerateID(),
		Name:     request.Name,
		Email:    request.Email,
		LogoUrl:  request.LogoUrl,
		Domain:   request.Domain,
		IsActive: 1,
	}

	if errCreateTenant := u.tenantRepo.Create(tenant); errCreateTenant != nil {
		log.WithFields(log.Fields{
			"error": errCreateTenant,
		}).Error("failed to create tenant")
		return fmt.Errorf("something Went wrong")
	}

	subscription := &modelTenant.Subscription{
		ID:        utils.GenerateID(),
		TenantID:  tenant.ID,
		PlanID:    request.SubscriptionPlanID,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, subscriptionPlan.DurationDay),
	}

	if errCreateSubscription := u.subscriptionRepo.Create(subscription); errCreateSubscription != nil {
		log.WithFields(log.Fields{
			"error": errCreateSubscription,
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
	if errCreateUser := u.userRepo.Create(user); errCreateUser != nil {
		log.WithFields(log.Fields{
			"error": errCreateUser,
		}).Error("failed to create user")
		return fmt.Errorf("something Went wrong")
	}

	return nil
}
