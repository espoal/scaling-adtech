package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/service"
)

type AdHandler struct {
	adService *service.AdService
	log       *zap.SugaredLogger
}

func NewAdHandler(adService *service.AdService, log *zap.SugaredLogger) *AdHandler {
	return &AdHandler{
		adService: adService,
		log:       log,
	}
}

func (h *AdHandler) GetWinningAds(c *fiber.Ctx) error {

	placement := c.Query("placement")
	category := c.Query("category")
	keyword := c.Query("keyword")

	ads, err := h.adService.GetWinningAds(placement, category, keyword)
	if err != nil {
		h.log.Errorw("Failed to find matching line items", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to find matching line items",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"winning_ads": ads,
	})
}
