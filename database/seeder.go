package database

import (
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

// SeedUsers inserts initial data into the users table
func SeedUsers(DB *gorm.DB) {

	var countRole, countSubscriptionPlan int64

	DB.Model(&modelUser.Role{}).Count(&countRole)
	DB.Model(&modelTenant.SubscriptionPlan{}).Count(&countSubscriptionPlan)
	// Pastikan DB sudah terkoneksi

	db := GetConnection()

	roles := []modelUser.Role{
		{
			ID:          uuid.New().String(),
			Name:        "Admin",
			Description: "Orang yang mengatur Tenant super Admin",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Organizer",
			Description: "Orang yang mengatur Event ",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Attendee",
			Description: "Orang yang menghadiri event ",
		},
	}

	subscriptionPlans := []modelTenant.SubscriptionPlan{
		{
			ID:          uuid.New().String(),
			Name:        "Free",
			Price:       0,
			DurationDay: 30,
		},
		{
			ID:          uuid.New().String(),
			Name:        "Basic",
			Price:       5000,
			DurationDay: 30,
		},
		{
			ID:          uuid.New().String(),
			Name:        "Premium",
			Price:       20000,
			DurationDay: 30,
		},
	}

	for _, tenant := range roles {
		if countRole == 0 {
			db.Create(&tenant)
		}
	}

	for _, subscriptionPlan := range subscriptionPlans {
		if countSubscriptionPlan == 0 {
			db.Create(&subscriptionPlan)
		}
	}

	log.Println("Tenant seeding completed")

}
