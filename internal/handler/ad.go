package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AdHandler struct {
	log *zap.SugaredLogger
}

func NewAdHandler(log *zap.SugaredLogger) *AdHandler {
	return &AdHandler{
		log: log,
	}
}

func (h *AdHandler) GetWinningAds(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Winning line item now",
	})
}
