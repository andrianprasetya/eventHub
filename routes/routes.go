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
	v1ApiPublic := c.Group("/api/v1")

	v1ApiPublic.Post("/issueToken", authHandler.Login)
	v1ApiPublic.Post("/register-tenant", tenantHandler.RegisterTenant)

	subscriptionPublic := v1ApiPublic.Group("/subscription")
	subscriptionPublic.Get("/get-plan-all", subscriptionPlanHandler.GetAll)
	subscriptionPublic.Get("/get-plan/:id", subscriptionPlanHandler.Get)

	//with auth
	v1ApiPrivate := c.Group("/api/v1", middleware.AuthMiddleware(redisClient))

	domain := v1ApiPrivate.Group("/:domain")

	//tenant routes
	user := domain.Group("/user")
	tenant := domain.Group("/tenant")

	subscriptionPrivate := domain.Group("/subscription")

	user.Post("/create", userHandler.Create)

	//tenant
	tenant.Post("/update/:id", tenantHandler.UpdateTenant)

	//subscription
	subscriptionPrivate.Post("/create", subscriptionPlanHandler.Create, middleware.RequireRole("super-admin"))
	subscriptionPrivate.Post("/update/:id", subscriptionPlanHandler.Update, middleware.RequireRole("super-admin"))
	subscriptionPrivate.Delete("/delete/:id", subscriptionPlanHandler.Delete, middleware.RequireRole("super-admin"))

}
