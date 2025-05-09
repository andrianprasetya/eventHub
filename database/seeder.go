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
