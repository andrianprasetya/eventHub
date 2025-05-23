package helper

import (
	"context"
	"fmt"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	"github.com/andrianprasetya/eventHub/internal/audit_security_log/repository"
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/shared/response"
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

func FilterEventCreate(req request.CreateEventRequest, errorMessages []response.FieldErrors) []response.FieldErrors {
	if req.EndDate.Before(req.StartDate) {
		errorMessages = append(errorMessages, response.FieldErrors{
			Field:   "end_date",
			Message: "end_date must not be before start_date",
		})
	}

	for i, discount := range req.Discounts {
		if discount.EndDate.Before(discount.StartDate) {
			errorMessages = append(errorMessages, response.FieldErrors{
				Field:   fmt.Sprintf("discounts[%d].end_date", i),
				Message: "end_date must not be before start_date",
			})
		}
		if discount.StartDate.Before(req.StartDate) {
			errorMessages = append(errorMessages, response.FieldErrors{
				Field:   fmt.Sprintf("discounts[%d].start_date", i),
				Message: "start_date discount must not be before start_date event",
			})
		}
		if discount.EndDate.After(req.EndDate) {
			errorMessages = append(errorMessages, response.FieldErrors{
				Field:   fmt.Sprintf("discounts[%d].end_date", i),
				Message: "end_date discount must not be after end_date event",
			})
		}
	}

	for i, session := range req.Sessions {
		if session.EndDateTime.Before(session.StartDateTime) {
			errorMessages = append(errorMessages, response.FieldErrors{
				Field:   fmt.Sprintf("sessions[%d].end_date_time", i),
				Message: "end_date_time session must not be before start_date_time",
			})
		}
		if session.StartDateTime.Before(req.StartDate) {
			errorMessages = append(errorMessages, response.FieldErrors{
				Field:   fmt.Sprintf("sessions[%d].start_date_time", i),
				Message: "start_date_time session must not be before event start_date",
			})
		}
		if session.EndDateTime.After(req.EndDate) {
			errorMessages = append(errorMessages, response.FieldErrors{
				Field:   fmt.Sprintf("sessions[%d].end_date_time", i),
				Message: "end_date_time session must not be after event end_date",
			})
		}
	}

	return errorMessages
}
