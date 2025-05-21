package main

import (
	"flag"
	"fmt"
	_ "github.com/andrianprasetya/eventHub/database/dialect/postgres"
	logRepository "github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	eventRepository "github.com/andrianprasetya/eventHub/internal/event/repository"
	eventUsecase "github.com/andrianprasetya/eventHub/internal/event/usecase"
	_ "github.com/andrianprasetya/eventHub/internal/shared/config"
	repositoryShared "github.com/andrianprasetya/eventHub/internal/shared/repository"
	tenantRepository "github.com/andrianprasetya/eventHub/internal/tenant/repository"
	tenantUsecase "github.com/andrianprasetya/eventHub/internal/tenant/usecase"
	ticketRepository "github.com/andrianprasetya/eventHub/internal/ticket/repository"
	userRepository "github.com/andrianprasetya/eventHub/internal/user/repository"
	userUsecase "github.com/andrianprasetya/eventHub/internal/user/usecase"
	"github.com/andrianprasetya/eventHub/routes"
	appServer "github.com/andrianprasetya/eventHub/server"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	//Argument parser for non-sensitive configurations
	host := flag.String("host", os.Getenv("APP_HOST"), "API server host")
	port := flag.String("port", os.Getenv("APP_PORT"), "API server port")
	flag.Parse()

	//initialize fiber
	app := fiber.New()

	//connectDB
	db := appServer.InitDatabase()
	redis := appServer.InitRedis()

	//Repository
	txManager := &repositoryShared.GormTxManager{DB: db}
	tenantRepo := tenantRepository.NewTenantRepository(db)
	tenantSettingRepo := tenantRepository.NewTenantSettingRepository(db)
	subscriptionRepo := tenantRepository.NewSubscriptionRepository(db)
	subscriptionPlanRepo := tenantRepository.NewSubscriptionPlanRepository(db)
	userRepo := userRepository.NewUserRepository(db)
	roleRepo := userRepository.NewRoleRepository(db)
	eventCategoryRepo := eventRepository.NewEventCategoryRepository(db)
	eventTagRepo := eventRepository.NewEventTagRepository(db)
	eventRepo := eventRepository.NewEventRepository(db)
	eventSessionRepo := eventRepository.NewEventSessionRepository(db)
	ticketRepo := ticketRepository.NewTicketRepository(db)
	discountRepo := ticketRepository.NewDiscountRepository(db)
	loginHistoryRepo := logRepository.NewLoginHistoryRepository(db)
	logActivityRepo := logRepository.NewLogActivityRepository(db)

	//Usecase
	tenantUC := tenantUsecase.NewTenantUsecase(txManager,
		tenantRepo,
		tenantSettingRepo,
		subscriptionRepo,
		subscriptionPlanRepo,
		userRepo,
		roleRepo,
		eventTagRepo,
		eventCategoryRepo)
	subscriptionPlanUC := tenantUsecase.NewSubscriptionPlanUsecase(subscriptionPlanRepo)
	userUC := userUsecase.NewUserUsecase(txManager, userRepo, roleRepo, tenantSettingRepo, loginHistoryRepo, logActivityRepo)
	roleUC := userUsecase.NewRoleUsecase(roleRepo)
	eventUC := eventUsecase.NewEventUsecase(
		txManager,
		tenantSettingRepo,
		eventRepo,
		eventTagRepo,
		eventCategoryRepo,
		eventSessionRepo,
		ticketRepo,
		discountRepo,
		logActivityRepo)

	routes.SetupRoutes(app, redis, tenantUC, subscriptionPlanUC, userUC, roleUC, eventUC)

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", *host, *port)))
}
