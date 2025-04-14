package service

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"sweng-task/internal/model"
)

type TrackingService struct {
	log       *zap.SugaredLogger
	validator *validator.Validate
}

func NewTrackingService(log *zap.SugaredLogger) *TrackingService {
	return &TrackingService{
		log: log,
	}
}

func (s *TrackingService) TrackEvent(event model.TrackingEvent) error {
	s.log.Infow("Tracking event",
		"event_type", event.EventType,
		"line_item_id", event.LineItemID,
	)

	// Here you would implement the logic to track the event
	// For example, you might want to store it in a database or send it to an analytics service

	err := publishToQueue("tracking_events", event)
	if err != nil {
		s.log.Errorw("Failed to publish event to queue",
			"error", err,
			"event_type", event.EventType,
			"line_item_id", event.LineItemID,
		)
		return err
	}

	return nil
}

func publishToQueue(queue string, event interface{}) error {

	return nil
}
