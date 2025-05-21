package service

import (
	"errors"
	request2 "github.com/andrianprasetya/eventHub/internal/event/dto/request"
	modelEvent "github.com/andrianprasetya/eventHub/internal/event/model"
	"github.com/andrianprasetya/eventHub/internal/shared/middleware"
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
)

var categories = []string{
	"Conference",
	"Workshop",
	"Seminar",
	"Webinar",
	"Concert",
	"Exhibition",
	"Meetup",
	"Networking",
	"Training",
	"Festival",
	"Product Launch",
	"Job Fair",
	"Fundraising",
	"Competition",
	"Religious",
	"Sports",
	"Ceremony",
	"Private Event",
	"Panel Discussion",
	"Hackathon",
}

var tags = []string{
	"Technology",
	"Marketing",
	"Education",
	"Finance",
	"Health",
	"Business",
	"Startup",
	"Environment",
	"Design",
	"AI/ML",
	"Crypto",
	"Photography",
}

func MapEventPayload(authUser interface{}, request interface{}) (*modelEvent.Event, error) {
	//type assertion if auth interface have a value AuthUser is true
	auth, ok := authUser.(middleware.AuthUser)
	if !ok {
		return nil, errors.New("mapping event payload failed")
	}
	req, ok := request.(request2.CreateEventRequest)
	if !ok {
		return nil, errors.New("mapping event payload failed")
	}

	return &modelEvent.Event{
		ID:          utils.GenerateID(),
		Title:       req.Title,
		TenantID:    auth.Tenant.ID,
		CategoryID:  req.CategoryID,
		EventType:   req.EventType,
		Tags:        req.Tags,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedBy:   auth.ID,
		IsTicket:    req.IsTicket,
		Status:      req.Status,
	}, nil
}

func BulkCategories(tenantID string) []*modelEvent.EventCategory {
	eventCategories := make([]*modelEvent.EventCategory, 0, len(categories))
	for _, t := range categories {
		eventCategories = append(eventCategories, &modelEvent.EventCategory{
			ID:       utils.GenerateID(),
			TenantID: tenantID,
			Name:     t,
		})
	}
	return eventCategories
}

func BulkTags(tenantID string) []*modelEvent.EventTag {
	eventTags := make([]*modelEvent.EventTag, 0, len(tags))
	for _, t := range tags {
		eventTags = append(eventTags, &modelEvent.EventTag{
			ID:       utils.GenerateID(),
			TenantID: tenantID,
			Name:     t,
		})
	}
	return eventTags
}
