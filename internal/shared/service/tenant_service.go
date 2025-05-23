package service

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/shared/validation"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
)

func MapTenantPayload(request request.CreateTenantRequest) *modelTenant.Tenant {
	return &modelTenant.Tenant{
		ID:       utils.GenerateID(),
		Name:     request.Name,
		Email:    request.Email,
		LogoUrl:  request.LogoUrl,
		Domain:   utils.GenerateDomainName(request.Name),
		IsActive: 1,
	}
}

func CheckMaxEventCanCreated(countEventCreated int, maxEvent int) error {
	if countEventCreated >= maxEvent {
		return validation.ValidationError{
			"event_limit": fmt.Sprintf("your subscription package not able create event more than %d", maxEvent),
		}
	}
	return nil
}
