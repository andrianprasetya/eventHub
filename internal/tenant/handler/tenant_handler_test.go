package handler

import (
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type MockTenantUsecase struct {
	mock.Mock
}

func (m *MockTenantUsecase) RegisterTenant(req request.CreateTenantRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockTenantUsecase) Update(id string, req request.UpdateTenantRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func TestRegisterTenant_Success(t *testing.T) {

	os.Setenv("MODE", "test")
	defer os.Unsetenv("MODE")

	app := fiber.New()
	mockUC := new(MockTenantUsecase)
	handler := NewTenantHandler(mockUC)

	body := `{"name":"Tenant Test","email":"tenant@test.com","logo_url": "test1234","password":"test1234","subscription_plan_id":"6bc51c92-d9b5-4c28-b41d-17092fa3c469"}`
	var reqParsed request.CreateTenantRequest
	_ = app.Config().JSONDecoder([]byte(body), &reqParsed)

	mockUC.On("RegisterTenant", reqParsed).Return(nil)
	app.Post("/register", handler.RegisterTenant)
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
	mockUC.AssertExpectations(t)
}
