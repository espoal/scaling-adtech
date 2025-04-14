package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/service"
)

type TrackingHandler struct {
	trackingService *service.TrackingService
	log             *zap.SugaredLogger
}

func NewTrackingHandler(trackingService *service.TrackingService, log *zap.SugaredLogger) *TrackingHandler {
	return &TrackingHandler{
		trackingService: trackingService,
		log:             log,
	}
}

func (h *TrackingHandler) TrackEvent(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
