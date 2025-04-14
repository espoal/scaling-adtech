package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/model"
	"sweng-task/internal/service"
)

type TrackingHandler struct {
	trackingService *service.TrackingService
	validator       *validator.Validate
	log             *zap.SugaredLogger
}

func NewTrackingHandler(trackingService *service.TrackingService, validator *validator.Validate, log *zap.SugaredLogger) *TrackingHandler {
	return &TrackingHandler{
		trackingService: trackingService,
		validator:       validator,
		log:             log,
	}
}

func (h *TrackingHandler) TrackEvent(c *fiber.Ctx) error {

	var input model.TrackingEvent
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := h.validator.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Validation error",
			"details": err.Error(),
		})
	}

	if err := h.trackingService.TrackEvent(input); err != nil {
		h.log.Errorw("Failed to track event", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to track event",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
