package handler

import (
	"log"
	"net/http"

	"chatbot/internal/models"
	"chatbot/internal/service"

	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	chatSvc service.ChatService
}

func NewChatHandler(chatSvc service.ChatService) *ChatHandler {
	return &ChatHandler{chatSvc: chatSvc}
}

func (h *ChatHandler) HandleChat(c echo.Context) error {
	var req models.ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Message == "" || req.SessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "message and session_id are required"})
	}

	reply, err := h.chatSvc.ProcessMessage(c.Request().Context(), req)
	if err != nil {
		log.Printf("chat error session=%s: %v", req.SessionID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to process message"})
	}

	return c.JSON(http.StatusOK, models.ChatResponse{Reply: reply})
}

func (h *ChatHandler) HandleGetSessions(c echo.Context) error {
	sessions, err := h.chatSvc.ListSessions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch sessions"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"sessions": sessions})
}

func (h *ChatHandler) HandleGetChat(c echo.Context) error {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "session_id is required"})
	}

	msgs, err := h.chatSvc.GetHistory(c.Request().Context(), sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch history"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"session_id": sessionID,
		"messages":   msgs,
	})
}
