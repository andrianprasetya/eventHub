package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/response"
	"github.com/andrianprasetya/eventHub/internal/event/model"
)

func FromEventModel(event *model.Event) *response.EventResponse {
	return &response.EventResponse{
		ID:    event.ID,
		Title: event.Title,
		Category: response.Category{
			ID:   event.Category.ID,
			Name: event.Category.Name,
		},
		Tags:        *event.Tags,
		Description: *event.Description,
		Location:    event.Location,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Status:      event.Status,
	}
}

func FromEventToListItem(event *model.Event) *response.EventListItemResponse {
	return &response.EventListItemResponse{
		ID:           event.ID,
		Title:        event.Title,
		CategoryName: event.Category.Name,
		Tags:         *event.Tags,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		Status:       event.Status,
	}
}

func FromEventToList(events []*model.Event) []*response.EventListItemResponse {
	result := make([]*response.EventListItemResponse, 0, len(events))
	for _, event := range events {
		result = append(result, FromEventToListItem(event))
	}
	return result
}

func FromEventTagToListItem(eventTag *model.EventTag) *response.EventTagListItemResponse {
	return &response.EventTagListItemResponse{
		Name: eventTag.Name,
	}
}

func FromEventTagToList(eventTags []*model.EventTag) []*response.EventTagListItemResponse {
	result := make([]*response.EventTagListItemResponse, 0, len(eventTags))
	for _, eventTag := range eventTags {
		result = append(result, FromEventTagToListItem(eventTag))
	}
	return result
}

func FromEventCategoryToListItem(eventCategory *model.EventCategory) *response.EventCategoryListItemResponse {
	return &response.EventCategoryListItemResponse{
		ID:   eventCategory.ID,
		Name: eventCategory.Name,
	}
}

func FromEventCategoryToList(eventCategories []*model.EventCategory) []*response.EventCategoryListItemResponse {
	result := make([]*response.EventCategoryListItemResponse, 0, len(eventCategories))
	for _, eventCategory := range eventCategories {
		result = append(result, FromEventCategoryToListItem(eventCategory))
	}
	return result
}
