package service

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
)

func MapTenantSettingsPayload(features map[string]interface{}, tenantID string) []*modelTenant.TenantSetting {
	var tenantSettings []*modelTenant.TenantSetting
	for key, value := range features {
		strVal := fmt.Sprintf("%v", value)
		tenantSettings = append(tenantSettings, &modelTenant.TenantSetting{
			ID:       utils.GenerateID(),
			TenantID: tenantID,
			Key:      key,
			Value:    strVal,
		})
	}
	return tenantSettings
}
