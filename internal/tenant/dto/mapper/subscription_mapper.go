package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/response"
	"github.com/andrianprasetya/eventHub/internal/tenant/model"
	log "github.com/sirupsen/logrus"
)

/*func FromUserModel(tenant *model2.Tenant) *response2.TenantResponse {
	return &response2.TenantResponse{}
}*/

func FromSubscriptionPlanToListItem(subscriptionPlan *model.SubscriptionPlan) *response.SubscriptionPlanListItemResponse {
	feature, err := utils.ToStringJSON(subscriptionPlan.Feature)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to un-marshal string feature")
	}

	return &response.SubscriptionPlanListItemResponse{
		ID:          subscriptionPlan.ID,
		Name:        subscriptionPlan.Name,
		Feature:     feature,
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
