package handler

import (
	"net/http"
	"time"

	"chatbot/internal/models"
	"chatbot/internal/repository"

	"github.com/labstack/echo/v4"
)

type LogHandler struct {
	logRepo repository.LogRepository
}

func NewLogHandler(logRepo repository.LogRepository) *LogHandler {
	return &LogHandler{logRepo: logRepo}
}

func (h *LogHandler) HandleLogs(c echo.Context) error {
	var log models.InferenceLog
	if err := c.Bind(&log); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid log payload"})
	}

	if log.Model == "" || log.Provider == "" || log.SessionID == "" || log.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "model, provider, session_id and status are required"})
	}

	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now().UTC()
	}

	if err := h.logRepo.InsertLog(c.Request().Context(), log); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to store log"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}
