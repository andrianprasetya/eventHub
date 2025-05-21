package service

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/request"
	"github.com/andrianprasetya/eventHub/internal/event/model"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
)

func MapEventServicesPayload(eventID string, sessions []request.EventSession) []*model.EventSession {
	eventSessions := make([]*model.EventSession, 0, len(sessions))
	for _, session := range sessions {
		eventSessions = append(eventSessions, &model.EventSession{
			ID:            utils.GenerateID(),
			EventID:       eventID,
			Title:         session.Title,
			StartDateTime: session.StartDateTime,
			EndDateTime:   session.EndDateTime,
		})
	}
	return eventSessions
}
