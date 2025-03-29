package main

import (
	"github.com/andrianprasetya/eventHub/database"
	"github.com/andrianprasetya/eventHub/routes"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	// Init Echo
	e := echo.New()

	// Database setup
	_, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup routes
	routes.SetupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
