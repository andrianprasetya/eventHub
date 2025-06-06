package routes

import (
	EH "github.com/andrianprasetya/eventHub/internal/event/handler"
	eventUsecase "github.com/andrianprasetya/eventHub/internal/event/usecase"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/redisser"
	TH "github.com/andrianprasetya/eventHub/internal/tenant/handler"
	tenantUsecase "github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	ETH "github.com/andrianprasetya/eventHub/internal/ticket/handler"
	ticketUsecase "github.com/andrianprasetya/eventHub/internal/ticket/usecase"
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
	eventUC eventUsecase.EventUsecase,
	ticketUC ticketUsecase.TicketUsecase,
	discountUC ticketUsecase.DiscountUsecase) {

	tenantHandler := TH.NewTenantHandler(tenantUC)
	subscriptionPlanHandler := TH.NewSubscriptionPlanHandler(subscriptionPlanUC)
	authHandler := UH.NewAuthHandler(userUC)
	userHandler := UH.NewUserHandler(userUC)
	roleHandler := UH.NewRoleHandler(roleUC)
	eventHandler := EH.NewEventHandler(eventUC)
	ticketHandler := ETH.NewTicketHandler(ticketUC)
	discountHandler := ETH.NewDiscountHandler(discountUC)

	//without auth
	v1ApiPublic := c.Group("/api/v1")

	v1ApiPublic.Post("/issueToken", authHandler.Login)
	v1ApiPublic.Post("/register-tenant", tenantHandler.RegisterTenant)

	subscriptionPublic := v1ApiPublic.Group("/subscription")
	eventPublic := v1ApiPublic.Group("/event")
	eventPublic.Get("/get-all-guest", eventHandler.GetAll)
	eventPublic.Get("/get/:id", eventHandler.GetByID)
	subscriptionPublic.Get("/get-plan-all", subscriptionPlanHandler.GetAll)
	subscriptionPublic.Get("/get-plan/:id", subscriptionPlanHandler.Get)

	//with auth
	v1ApiPrivate := c.Group("/api/v1", middleware.AuthMiddleware(redisClient))

	user := v1ApiPrivate.Group("/user")
	role := v1ApiPrivate.Group("/role")
	tenant := v1ApiPrivate.Group("/tenant")
	event := v1ApiPrivate.Group("/event")
	ticket := v1ApiPublic.Group("/ticket")
	discount := v1ApiPublic.Group("/discount")
	subscription := v1ApiPrivate.Group("/subscription")

	//user
	user.Post("/create", middleware.RequireRole("tenant-admin"), userHandler.Create)
	user.Get("/get-all", middleware.RequireRole("tenant-admin", "tenant-admin"), userHandler.GetAll)
	user.Get("/get/:id", middleware.RequireRole("tenant-admin", "tenant-admin"), userHandler.GetByID)
	user.Post("/update/:id", middleware.RequireRole("tenant-admin"), userHandler.Update)
	user.Delete("/delete/:id", middleware.RequireRole("tenant-admin"), userHandler.Delete)

	//role
	role.Get("/get-all", middleware.RequireRole("tenant-admin", "super-admin"), roleHandler.GetAll)
	role.Get("/get/:id", middleware.RequireRole("tenant-admin", "super-admin"), roleHandler.GetByID)

	//tenant
	tenant.Post("/update-information/:id", middleware.RequireRole("tenant-admin"), tenantHandler.UpdateInformation)
	tenant.Post("/update-information/:id", middleware.RequireRole("tenant-admin"), tenantHandler.UpdateInformation)

	//event
	event.Get("/get-tags", middleware.RequireRole("tenant-admin", "organizer"), eventHandler.GetTags)
	event.Get("/get-categories", middleware.RequireRole("tenant-admin", "organizer"), eventHandler.GetCategories)
	event.Post("/create", middleware.RequireRole("tenant-admin", "organizer"), eventHandler.Create)
	event.Get("/get-all", middleware.RequireRole("super-admin", "tenant-admin", "organizer"), eventHandler.GetAll)

	//ticket
	ticket.Post("/create", ticketHandler.Create)
	ticket.Get("/get-all", ticketHandler.GetAll)
	ticket.Get("/get/:id", ticketHandler.GetByID)

	//discount
	discount.Post("/create", discountHandler.Create)
	discount.Get("/get-all", discountHandler.GetAll)
	discount.Get("/get/:id", discountHandler.GetByID)

	//subscription
	subscription.Post("/create", middleware.RequireRole("super-admin"), subscriptionPlanHandler.Create)
	subscription.Post("/update/:id", middleware.RequireRole("super-admin"), subscriptionPlanHandler.Update)
	subscription.Delete("/delete/:id", middleware.RequireRole("super-admin"), subscriptionPlanHandler.Delete)
}
