package service

import (
	"go.uber.org/zap"
)

type TrackingService struct {
	log             *zap.SugaredLogger
	lineItemService *LineItemService
}

func NewTrackingService(log *zap.SugaredLogger) *TrackingService {
	return &TrackingService{
		log: log,
	}
}
