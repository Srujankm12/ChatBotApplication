package models

import "time"

type Message struct {
	Role      string    `bson:"role" json:"role"`
	Content   string    `bson:"content" json:"content"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Chat struct {
	SessionID string    `bson:"session_id" json:"session_id"`
	Messages  []Message `bson:"messages" json:"messages"`
}

type ChatRequest struct {
	Message   string `json:"message" validate:"required"`
	SessionID string `json:"session_id" validate:"required"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}
