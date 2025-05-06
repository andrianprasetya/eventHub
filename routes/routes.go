package routes

import (
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/redisser"
	TH "github.com/andrianprasetya/eventHub/internal/tenant/handler"
	tenantUsecase "github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	UH "github.com/andrianprasetya/eventHub/internal/user/handler"
	userUsecase "github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(c *fiber.App,
	redisClient redisser.RedisClient,
	tenantUC tenantUsecase.TenantUsecase,
	subscriptionPlanUC tenantUsecase.SubscriptionPlanUsecase,
	userUC userUsecase.UserUsecase) {

	tenantHandler := TH.NewTenantHandler(tenantUC)
	subscriptionPlanHandler := TH.NewSubscriptionPlanHandler(subscriptionPlanUC)
	authHandler := UH.NewAuthHandler(userUC)
	userHandler := UH.NewUserHandler(userUC)
	api := c.Group("/api")
	auth := c.Group("/auth")
	v1 := api.Group("/v1")

	//tenant routes
	tenant := v1.Group("/tenant")
	subscription := v1.Group("/subscription")

	user := v1.Group("/user", middleware.AuthMiddleware(redisClient))

	auth.Post("/issueToken", authHandler.Login)
	user.Post("/create", userHandler.Create)
	tenant.Post("/register", tenantHandler.RegisterTenant)
	subscription.Get("/get-plan", subscriptionPlanHandler.GetAll)
}
