package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/event/dto/response"
	"github.com/andrianprasetya/eventHub/internal/event/model"
)

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
