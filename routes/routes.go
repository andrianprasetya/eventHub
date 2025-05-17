package routes

import (
	EH "github.com/andrianprasetya/eventHub/internal/event/handler"
	eventUsecase "github.com/andrianprasetya/eventHub/internal/event/usecase"
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
	userUC userUsecase.UserUsecase,
	roleUC userUsecase.RoleUsecase,
	eventUC eventUsecase.EventUsecase) {

	tenantHandler := TH.NewTenantHandler(tenantUC)
	subscriptionPlanHandler := TH.NewSubscriptionPlanHandler(subscriptionPlanUC)
	authHandler := UH.NewAuthHandler(userUC)
	userHandler := UH.NewUserHandler(userUC)
	roleHandler := UH.NewRoleHandler(roleUC)
	eventHandler := EH.NewEventHandler(eventUC)

	//without auth
	v1ApiPublic := c.Group("/api/v1")

	v1ApiPublic.Post("/issueToken", authHandler.Login)
	v1ApiPublic.Post("/register-tenant", tenantHandler.RegisterTenant)

	subscriptionPublic := v1ApiPublic.Group("/subscription")
	subscriptionPublic.Get("/get-plan-all", subscriptionPlanHandler.GetAll)
	subscriptionPublic.Get("/get-plan/:id", subscriptionPlanHandler.Get)

	//with auth
	v1ApiPrivate := c.Group("/api/v1", middleware.AuthMiddleware(redisClient))

	//group without domain
	subscriptionPrivate := v1ApiPrivate.Group("/subscription")
	userAdmin := v1ApiPrivate.Group("/user")

	//user
	userAdmin.Get("/get-all", userHandler.GetAll, middleware.RequireRole("super-admin"))
	userAdmin.Get("/get/:id", userHandler.GetByID, middleware.RequireRole("super-admin"))

	//subscription
	subscriptionPrivate.Post("/create", subscriptionPlanHandler.Create, middleware.RequireRole("super-admin"))
	subscriptionPrivate.Post("/update/:id", subscriptionPlanHandler.Update, middleware.RequireRole("super-admin"))
	subscriptionPrivate.Delete("/delete/:id", subscriptionPlanHandler.Delete, middleware.RequireRole("super-admin"))

	domain := v1ApiPrivate.Group("/:domain")

	//domain groups
	user := domain.Group("/user")
	role := domain.Group("/role")
	tenant := domain.Group("/tenant")
	event := domain.Group("/event")

	//user
	user.Post("/create", userHandler.Create, middleware.RequireRole("tenant-admin"))
	user.Get("/get-all", userHandler.GetAll)
	user.Get("/get/:id", userHandler.GetByID)
	user.Post("/update/:id", userHandler.Update, middleware.RequireRole("tenant-admin"))
	user.Delete("/delete/:id", userHandler.Delete, middleware.RequireRole("tenant-admin"))

	//role
	role.Get("/get-all", roleHandler.GetAll)
	role.Get("/get/:id", roleHandler.GetByID)

	//tenant
	tenant.Post("/update-information/:id", tenantHandler.UpdateInformation, middleware.RequireRole("tenant-admin"))
	tenant.Post("/update-information/:id", tenantHandler.UpdateInformation, middleware.RequireRole("tenant-admin"))

	//event
	event.Get("/get-tags", eventHandler.GetTags, middleware.RequireRole("tenant-admin", "organizer"))
	event.Get("/get-categories", eventHandler.GetCategories, middleware.RequireRole("tenant-admin", "organizer"))
	event.Post("/create", eventHandler.Create, middleware.RequireRole("tenant-admin", "organizer"))

}
