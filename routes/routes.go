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

	//without auth
	apiWithoutAuth := c.Group("/api")
	v1WithoutAuth := apiWithoutAuth.Group("/v1")

	//with auth
	apiWithAuth := c.Group("/api", middleware.AuthMiddleware(redisClient))
	v1WithAuth := apiWithAuth.Group("/v1")

	domain := v1WithAuth.Group("/:domain")

	//tenant routes
	user := domain.Group("/user")
	tenant := domain.Group("/tenant")
	subscriptionWithoutAuth := v1WithoutAuth.Group("/subscription")
	subscriptionWithAuth := v1WithAuth.Group("/subscription")

	apiWithoutAuth.Post("/issueToken", authHandler.Login)
	apiWithoutAuth.Post("/register-tenant", tenantHandler.RegisterTenant)

	user.Post("/create", userHandler.Create)

	//tenant
	tenant.Post("/update/:id", tenantHandler.UpdateTenant)

	//subscription
	subscriptionWithAuth.Post("create", subscriptionPlanHandler.Create, middleware.RequireRole("super-admin"))
	subscriptionWithoutAuth.Get("/get-plan", subscriptionPlanHandler.GetAll)
}
