package response

import "time"

type FieldErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type UserLog struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	RoleID   string `json:"role_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive int    `json:"is_active"`
}

type EventLog struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Description *string   `json:"description"`
	Location    string    `json:"location"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
