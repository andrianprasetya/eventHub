package mapper

import (
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/response"
	"github.com/andrianprasetya/eventHub/internal/ticket/model"
)

func FromDiscountModel(discount *model.Discount) *response.DiscountResponse {
	return &response.DiscountResponse{
		ID:                 discount.ID,
		EventID:            discount.EventID,
		Code:               discount.Code,
		DiscountPercentage: discount.DiscountPercentage,
		StartDate:          discount.StartDate,
		EndDate:            discount.EndDate,
	}
}

func FromDiscountToListItem(discount *model.Discount) *response.DiscountListItemResponse {
	return &response.DiscountListItemResponse{
		ID:        discount.ID,
		Code:      discount.Code,
		StartDate: discount.StartDate,
		EndDate:   discount.EndDate,
	}
}

func FromDiscountToList(discounts []*model.Discount) []*response.DiscountListItemResponse {
	result := make([]*response.DiscountListItemResponse, 0, len(discounts))
	for _, discount := range discounts {
		result = append(result, FromDiscountToListItem(discount))
	}
	return result
}
