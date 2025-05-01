package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
)

/*func FromUserModel(tenant *model2.Tenant) *response2.TenantResponse {
	return &response2.TenantResponse{}
}*/

func FromSubscriptionPlanToListItem(subscriptionPlan *model.SubscriptionPlan) *response.SubscriptionPlanListItemResponse {
	return &response.SubscriptionPlanListItemResponse{
		ID:          subscriptionPlan.ID,
		Name:        subscriptionPlan.Name,
		Price:       subscriptionPlan.Price,
		DurationDay: subscriptionPlan.DurationDay,
	}
}

func FromSubscriptionPlanToList(subscriptionPlans []*model.SubscriptionPlan) []*response.SubscriptionPlanListItemResponse {
	result := make([]*response.SubscriptionPlanListItemResponse, 0, len(subscriptionPlans))
	for _, subscriptionPlan := range subscriptionPlans {
		result = append(result, FromSubscriptionPlanToListItem(subscriptionPlan))
	}
	return result
}
