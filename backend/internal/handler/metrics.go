package handler

import (
	"net/http"

	"chatbot/internal/service"

	"github.com/labstack/echo/v4"
)

type MetricsHandler struct {
	metricsSvc service.MetricsService
}

func NewMetricsHandler(metricsSvc service.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsSvc: metricsSvc}
}

func (h *MetricsHandler) HandleMetrics(c echo.Context) error {
	metrics, err := h.metricsSvc.GetMetrics(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch metrics"})
	}
	return c.JSON(http.StatusOK, metrics)
}
