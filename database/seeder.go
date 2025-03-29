package database

import (
	models2 "github.com/andrianprasetya/eventHub/internal/models"
	"github.com/google/uuid"
	"log"
)

// SeedUsers inserts initial data into the users table
func SeedUsers() {

	// Pastikan DB sudah terkoneksi
	if DB == nil {
		log.Fatal("Database is not connected")
	}

	tenants := []models2.Tenant{
		{
			ID:    uuid.New().String(),
			Name:  "Tenant Alpha",
			Email: "alpha@gmail.com",
		},
		{
			ID:    uuid.New().String(),
			Name:  "Tenant Beta",
			Email: "beta@gmail.com",
		},
	}

	for _, tenant := range tenants {
		DB.Create(&tenant)

		settings := []models2.TenantSetting{
			{ID: uuid.New().String(), TenantID: tenant.ID, Key: "timezone", Value: "UTC"},
			{ID: uuid.New().String(), TenantID: tenant.ID, Key: "currency", Value: "USD"},
			{ID: uuid.New().String(), TenantID: tenant.ID, Key: "language", Value: "en"},
			{ID: uuid.New().String(), TenantID: tenant.ID, Key: "theme", Value: "dark"},
		}

		for _, setting := range settings {
			DB.Create(&setting)
		}

	}
	log.Println("User seeding completed")
}
