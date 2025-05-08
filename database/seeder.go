package database

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
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
			ID:          utils.GenerateID(),
			Name:        "Admin",
			Slug:        utils.Slugify("Admin"),
			Description: "Orang yang mengatur Tenant super Admin",
		},
		{
			ID:          utils.GenerateID(),
			Name:        "Organizer",
			Slug:        utils.Slugify("Organizer"),
			Description: "Orang yang mengatur Event ",
		},
	}

	subscriptionPlans := []modelTenant.SubscriptionPlan{
		{
			ID:          utils.GenerateID(),
			Name:        "Free",
			Price:       0,
			DurationDay: 30,
		},
		{
			ID:          utils.GenerateID(),
			Name:        "Basic",
			Price:       5000,
			DurationDay: 30,
		},
		{
			ID:          utils.GenerateID(),
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
