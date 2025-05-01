package routes

import (
	"github.com/andrianprasetya/eventHub/internal/tenant"
	TH "github.com/andrianprasetya/eventHub/internal/tenant/handler"
	"github.com/andrianprasetya/eventHub/internal/user"
	UH "github.com/andrianprasetya/eventHub/internal/user/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(c *fiber.App, tenantUC tenant.TenantUsecase, subscriptionPlanUC tenant.SubscriptionPlanUsecase, userUC user.UserUsecase) {
	tenantHandler := TH.NewTenantHandler(tenantUC)
	subscriptionPlanHandler := TH.NewSubscriptionPlanHandler(subscriptionPlanUC)
	authHandler := UH.NewAuthHandler(userUC)
	api := c.Group("/api")
	v1 := api.Group("/v1")

	//tenant routes
	tenant := v1.Group("/tenant")
	subscription := v1.Group("/subscription")
	user := v1.Group("/user")

	user.Post("/login", authHandler.Login)
	tenant.Post("/register", tenantHandler.RegisterTenant)
	subscription.Get("/get-plan", subscriptionPlanHandler.GetAll)
}
