package database

import (
	models2 "github.com/andrianprasetya/eventHub/internal/models"
	"log"
)

// MigrateDatabase runs database migrations
func MigrateDatabase() {
	err := DB.AutoMigrate(
		&models2.Tenant{},
		&models2.TenantSetting{},
		&models2.User{},
		&models2.PasswordReset{},
		&models2.Event{},
		&models2.Ticket{},
		&models2.EventTicket{},
		&models2.Order{},
		&models2.OrderItem{},
		&models2.CheckIn{},
		&models2.ActivityLog{},
		&models2.Notification{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration completed successfully")
}
