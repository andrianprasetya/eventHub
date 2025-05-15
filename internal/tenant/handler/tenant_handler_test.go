package handler

import (
	"encoding/json"
	"github.com/andrianprasetya/eventHub/internal/tenant/dto/request"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestRegisterTenant_Failed(t *testing.T) {
	os.Setenv("MODE", "test")
	defer os.Unsetenv("MODE")

	app := fiber.New()
	mockUC := new(MockTenantUsecase)
	handler := NewTenantHandler(mockUC)

	app.Post("/register", handler.RegisterTenant)

	// Password hanya 4 karakter: gagal validasi min=8
	reqBody := request.CreateTenantRequest{
		Name:               "Tenant Test",
		Email:              "tenant@test.com",
		LogoUrl:            "test1234",
		Password:           "test1234",
		SubscriptionPlanID: "6bc51c92-d9b5-4c28-b41d-17092fa3c469",
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	body := string(jsonBytes)

	mockUC.On("RegisterTenant", reqBody).Return(nil)
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Cetak response body untuk debug
	/*resBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Validation Error Response:", string(resBody))

	// Kamu juga bisa unmarshal dan assert isi response error
	var result map[string]interface{}
	json.Unmarshal(resBody, &result)

	assert.Contains(t, result, "errors")
	errorsMap := result["errors"].(map[string]interface{})
	assert.Contains(t, errorsMap, "password")*/
}
