package service

import (
	"context"

	"chatbot/internal/models"
	"chatbot/internal/repository"
)

type MetricsService interface {
	GetMetrics(ctx context.Context) (models.MetricsResponse, error)
}

type metricsService struct {
	logRepo repository.LogRepository
}

func NewMetricsService(logRepo repository.LogRepository) MetricsService {
	return &metricsService{logRepo: logRepo}
}

func (s *metricsService) GetMetrics(ctx context.Context) (models.MetricsResponse, error) {
	return s.logRepo.GetMetrics(ctx)
}
