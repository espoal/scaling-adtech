package service

import (
	"cmp"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"slices"
	"sweng-task/internal/model"
)

type AdService struct {
	log             *zap.SugaredLogger
	lineItemService *LineItemService
}

func NewAdService(log *zap.SugaredLogger, lineItemService *LineItemService) *AdService {
	return &AdService{
		log:             log,
		lineItemService: lineItemService,
	}
}

func (h *AdService) GetWinningAds(placement, category, keyword string) ([]*model.LineItem, error) {

	ads, err := h.lineItemService.FindMatchingLineItems(placement, category, keyword)
	if err != nil {
		h.log.Errorw("Failed to find matching line items", "error", err)
		return nil, err
	}

	// Relax constraints if no ads are found
	if len(ads) == 0 && category != "" {
		ads, err = h.lineItemService.GetAll("", placement, category, "", "active")
		if err != nil {
			log.Errorw("Failed to retrieve line items",
				"error", err,
				"placement", placement,
			)
			return nil, err
		}
	}

	slices.SortFunc(ads, func(a, b *model.LineItem) int {
		return cmp.Compare(a.Bid, b.Bid)
	})

	return ads, nil
}
