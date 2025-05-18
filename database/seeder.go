package database

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
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
			Name:        "Super Admin",
			Slug:        utils.Slugify("Super Admin"),
			Description: "Orang yang mengatur Segala yg berhubungan dengan aplikasi diluar kendali tenant",
		},
		{
			ID:          utils.GenerateID(),
			Name:        "Tenant Admin",
			Slug:        utils.Slugify("Tenant Admin"),
			Description: "Orang yang mengatur Tenant",
		},
		{
			ID:          utils.GenerateID(),
			Name:        "Organizer",
			Slug:        utils.Slugify("Organizer"),
			Description: "Orang yang mengatur Event ",
		},
	}

	featureFree, err := utils.ToJSONString(map[string]interface{}{
		"max_events":                  1,
		"max_tickets_per_event":       50,
		"max_users":                   1,
		"unlimited_events":            false,
		"unlimited_tickets_per_event": false,
		"unlimited_users":             false,
	})
	featureBasic, err := utils.ToJSONString(map[string]interface{}{
		"max_events":                  4,
		"max_tickets_per_event":       200,
		"max_users":                   4,
		"unlimited_events":            false,
		"unlimited_tickets_per_event": false,
		"unlimited_users":             false,
	})
	featurePremium, err := utils.ToJSONString(map[string]interface{}{
		"max_events":                  0,
		"max_tickets_per_event":       0,
		"max_users":                   0,
		"unlimited_events":            true,
		"unlimited_tickets_per_event": true,
		"unlimited_users":             true,
	})

	if err != nil {
		panic(err)
	}

	subscriptionPlans := []modelTenant.SubscriptionPlan{
		{
			ID:          utils.GenerateID(),
			Name:        "Free",
			Feature:     featureFree,
			DurationDay: -1,
		},
		{
			ID:          utils.GenerateID(),
			Name:        "Basic",
			Feature:     featureBasic,
			Price:       5000,
			DurationDay: 30,
		},
		{
			ID:          utils.GenerateID(),
			Name:        "Premium",
			Feature:     featurePremium,
			Price:       20000,
			DurationDay: 30,
		},
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SUPERADMIN_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error seeding bcrypt failed")
	}
	for _, role := range roles {
		if countRole == 0 {
			db.Create(&role)
			if role.Slug == "super-admin" {
				err := db.Create(&modelUser.User{
					ID:       utils.GenerateID(),
					TenantID: "",
					Name:     "Super Admin",
					Email:    "superadmin@gmail.com",
					RoleID:   role.ID,
					Password: string(hashedPassword),
					IsActive: 1,
				}).Error

				if err != nil {
					log.Fatal("Gagal insert user:", err)
				}
			}
		}
	}

	for _, subscriptionPlan := range subscriptionPlans {
		if countSubscriptionPlan == 0 {
			db.Create(&subscriptionPlan)
		}
	}

	log.Println("Tenant seeding completed")

}
