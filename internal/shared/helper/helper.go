package helper

import (
	"context"
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	logging "github.com/sirupsen/logrus"
	"time"
)

type CountEventResult struct {
	Count int
	Err   error
}

type UnlimitedResult struct {
	Value string
	Err   error
}

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
				"errors": err,
			}).Error("failed to Log Login History")
		}
	}(log)
}

func LogActivity(repo repository.LogActivityRepository, ctx context.Context, tenantID, userID, url, action, objectData, objectType, objectID string) {
	activity := &model.ActivityLog{
		ID:         utils.GenerateID(),
		TenantID:   tenantID,
		UserID:     userID,
		URL:        url,
		Action:     action,
		ObjectData: objectData,
		ObjectType: objectType,
		ObjectID:   objectID,
	}

	go func(activity *model.ActivityLog) {
		defer func() {
			if r := recover(); r != nil {
				logging.WithFields(logging.Fields{
					"recover": r,
				}).Error("panic occurred in LogLoginHistory goroutine")
			}
		}()
		if err := repo.Create(ctx, activity); err != nil {
			logging.WithFields(logging.Fields{
				"errors": err,
			}).Error("failed to Log Activity")
		}
	}(activity)
}

func FilterEventCreate(req request.CreateEventRequest, errorMessages map[string]string) map[string]string {
	if req.EndDate.Before(req.StartDate) {
		errorMessages["end_date"] = "end_date must not below start_date"
	}
	for i, discount := range req.Discounts {
		if discount.EndDate.Before(discount.StartDate) {
			key := fmt.Sprintf("discounts[%d].end_date", i)
			errorMessages[key] = "end_date must not below start_date"
		}
		if discount.StartDate.Before(req.StartDate) {
			key := fmt.Sprintf("discounts[%d].start_date", i)
			errorMessages[key] = "start_date discount must not below start_date event"
		}
		if discount.EndDate.After(req.EndDate) {
			key := fmt.Sprintf("discounts[%d].end_date", i)
			errorMessages[key] = "end_date discount must not above end_start event"
		}
	}
	for i, session := range req.Sessions {
		if session.EndDateTime.Before(session.StartDateTime) {
			key := fmt.Sprintf("sessions[%d].end_date_time", i)
			errorMessages[key] = "end_date_time sessions must not below start_date_time"
		}
		if session.StartDateTime.Before(req.StartDate) {
			key := fmt.Sprintf("sessions[%d].start_date_time", i)
			errorMessages[key] = "start_date sessions must not below start_date event"
		}
		if session.EndDateTime.After(req.EndDate) {
			key := fmt.Sprintf("sessions[%d].end_date_time", i)
			errorMessages[key] = "end_date sessions must not above end_start event"
		}
	}
	return errorMessages
}
