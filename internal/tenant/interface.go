package tenant

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
)

type TenantRepository interface {
	Create(tenant *modelTenant.Tenant) error
}

type SubscriptionPlanRepository interface {
	GetAll() ([]*modelTenant.SubscriptionPlan, error)
	Get(id string) (*modelTenant.SubscriptionPlan, error)
}

type SubscriptionRepository interface {
	Create(subscription *modelTenant.Subscription) error
}

type TenantUsecase interface {
	RegisterTenant(request request.CreateTenantRequest) error
}

type SubscriptionPlanUsecase interface {
	GetAll() ([]*response.SubscriptionPlanListItemResponse, error)
}
