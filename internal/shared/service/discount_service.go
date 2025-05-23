package service

import (
	"github.com/andrianprasetya/eventHub/internal/shared/utils"
	"github.com/andrianprasetya/eventHub/internal/ticket/dto/request"
	modelTicket "github.com/andrianprasetya/eventHub/internal/ticket/model"
)

func MapDiscountsPayload(eventID string, discounts []request.EventDiscount) []*modelTicket.Discount {
	eventDiscounts := make([]*modelTicket.Discount, 0, len(discounts))
	for _, discount := range discounts {
		eventDiscounts = append(eventDiscounts, &modelTicket.Discount{
			ID:                 utils.GenerateID(),
			EventID:            eventID,
			Code:               discount.Code,
			DiscountPercentage: discount.DiscountPercentage,
			StartDate:          discount.StartDate,
			EndDate:            discount.EndDate,
		})
	}
	return eventDiscounts
}
