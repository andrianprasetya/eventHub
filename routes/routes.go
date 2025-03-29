package routes

import (
	"github.com/andrianprasetya/eventHub/internal/handlers"
	"github.com/andrianprasetya/eventHub/internal/usecases"
	"github.com/andrianprasetya/eventHub/middleware"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, eventUC *usecases.EventUsecase) {
	// Middleware JWT untuk route yang membutuhkan autentikasi
	authMiddleware := middleware.JWTMiddleware("your-secret-key")
	// Buat instance handler dengan dependency injection
	eventHandler := handlers.NewEventHandler(eventUC)

	// Grouping routes
	api := e.Group("/api")

	// Auth routes
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	// Event routes (protected)
	events := api.Group("/events", authMiddleware)
	events.POST("", eventHandler.CreateEvent)
	api.GET("/events/:id", eventHandler.GetEventByID)

	// Ticket routes (protected)
	tickets := api.Group("/tickets", authMiddleware)
	tickets.POST("", handlers.CreateTicket)
	tickets.GET("/:id", handlers.GetTicket)
	tickets.GET("", handlers.GetAllTickets)
	tickets.DELETE("/:id", handlers.DeleteTicket)
}
