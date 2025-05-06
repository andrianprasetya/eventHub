package middleware

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
	logging "github.com/sirupsen/logrus"
	"time"
)

func AuthenticationMiddleware(repo repository.LoginHistoryRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		userLocal := ctx.Locals("user")
		user, errGetLocalUser := userLocal.(AuthUser)

		if !errGetLocalUser {
			return err
		}

		log := &model.LoginHistory{
			ID:        utils.GenerateID(),
			UserID:    user.ID,
			LoginTime: time.Now(),
			IpAddress: ctx.IP(),
		}
		fmt.Println("test")
		go func(log *model.LoginHistory) {
			if err := repo.Create(log); err != nil {
				logging.WithFields(logging.Fields{
					"error": err,
				}).Error("failed to Log Login History")
			}
		}(log)

		return err
	}
}
