package service

import (
	"context"
	"fmt"

	"chatbot/internal/models"
	"chatbot/internal/repository"
	"chatbot/sdk"
)

type ChatService interface {
	ProcessMessage(ctx context.Context, req models.ChatRequest) (string, error)
	GetHistory(ctx context.Context, sessionID string) ([]models.Message, error)
	ListSessions(ctx context.Context) ([]string, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{chatRepo: chatRepo}
}

func (s *chatService) ProcessMessage(ctx context.Context, req models.ChatRequest) (string, error) {
	if req.Message == "" {
		return "", fmt.Errorf("message cannot be empty")
	}
	if req.SessionID == "" {
		return "", fmt.Errorf("session_id cannot be empty")
	}

	history, err := s.chatRepo.GetMessages(ctx, req.SessionID)
	if err != nil {
		return "", fmt.Errorf("fetch history: %w", err)
	}

	userMsg := models.Message{Role: "user", Content: req.Message}
	history = append(history, userMsg)

	reply, err := sdk.CallLLM(ctx, history, req.SessionID)
	if err != nil {
		return "", fmt.Errorf("llm call: %w", err)
	}

	if err := s.chatRepo.AppendMessage(ctx, req.SessionID, userMsg); err != nil {
		return "", fmt.Errorf("save user message: %w", err)
	}

	assistantMsg := models.Message{Role: "assistant", Content: reply}
	if err := s.chatRepo.AppendMessage(ctx, req.SessionID, assistantMsg); err != nil {
		return "", fmt.Errorf("save assistant message: %w", err)
	}

	return reply, nil
}

func (s *chatService) GetHistory(ctx context.Context, sessionID string) ([]models.Message, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session_id cannot be empty")
	}
	return s.chatRepo.GetAllMessages(ctx, sessionID)
}

func (s *chatService) ListSessions(ctx context.Context) ([]string, error) {
	return s.chatRepo.ListSessions(ctx)
}
