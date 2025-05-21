package service

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"time"
)

func MapSubscriptionPayload(tenantID, planID string, durationDay int) *modelTenant.Subscription {
	var endDate *time.Time
	if durationDay != -1 {
		d := time.Now().AddDate(0, 0, durationDay)
		endDate = &d
	}
	return &modelTenant.Subscription{
		ID:        utils.GenerateID(),
		TenantID:  tenantID,
		PlanID:    planID,
		StartDate: time.Now(),
		EndDate:   endDate,
		IsActive:  1,
	}
}
