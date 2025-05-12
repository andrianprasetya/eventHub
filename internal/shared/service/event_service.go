package service

import (
	modelEvent "github.com/andrianprasetya/eventHub/internal/event/model"
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

func BulkCategories(tenantID string) *[]modelEvent.EventCategory {
	var eventCategories []modelEvent.EventCategory
	for _, t := range categories {
		eventCategories = append(eventCategories, modelEvent.EventCategory{
			ID:       utils.GenerateID(),
			TenantID: tenantID,
			Name:     t,
		})
	}
	return &eventCategories
}

func BulkTags(tenantID string) *[]modelEvent.EventTag {
	var eventTags []modelEvent.EventTag
	for _, t := range tags {
		eventTags = append(eventTags, modelEvent.EventTag{
			ID:       utils.GenerateID(),
			TenantID: tenantID,
			Name:     t,
		})
	}
	return &eventTags
}
