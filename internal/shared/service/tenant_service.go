package service

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	"strconv"
)

func CheckMaxEventCanCreated(countEventCreated int, tenantSetting *modelTenant.TenantSetting) error {
	maxEvent, _ := strconv.Atoi(tenantSetting.Value)

	if countEventCreated >= maxEvent {

		return validation.ValidationError{
			"event_limit": fmt.Sprintf("your subscription package not able create event more than %d", maxEvent),
		}
	}
	return nil
}
