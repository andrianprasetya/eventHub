package helper

import (
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	logging "github.com/sirupsen/logrus"
	"time"
)

func LogLoginHistory(repo repository.LoginHistoryRepository, userId, ip string) {

	log := &model.LoginHistory{
		ID:        utils.GenerateID(),
		UserID:    userId,
		LoginTime: time.Now(),
		IpAddress: ip,
	}
	go func(log *model.LoginHistory) {
		defer func() {
			if r := recover(); r != nil {
				logging.WithFields(logging.Fields{
					"recover": r,
				}).Error("panic occurred in LogLoginHistory goroutine")
			}
		}()

		if err := repo.Create(log); err != nil {
			logging.WithFields(logging.Fields{
				"error": err,
			}).Error("failed to Log Login History")
		}
	}(log)
}

func LogActivity(repo repository.LogActivityRepository, userId, url, action, objectData, objectType, objectId string) {
	activity := &model.ActivityLog{
		ID:         utils.GenerateID(),
		UserID:     userId,
		URL:        url,
		Action:     action,
		ObjectData: objectData,
		ObjectType: objectType,
		ObjectID:   objectId,
	}

	go func(activity *model.ActivityLog) {
		defer func() {
			if r := recover(); r != nil {
				logging.WithFields(logging.Fields{
					"recover": r,
				}).Error("panic occurred in LogLoginHistory goroutine")
			}
		}()
		if err := repo.Create(activity); err != nil {
			logging.WithFields(logging.Fields{
				"error": err,
			}).Error("failed to Log Activity")
		}
	}(activity)
}
