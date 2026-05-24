package models

import "time"

type InferenceLog struct {
	Model         string    `bson:"model" json:"model"`
	Provider      string    `bson:"provider" json:"provider"`
	LatencyMs     int64     `bson:"latency_ms" json:"latency_ms"`
	InputPreview  string    `bson:"input_preview" json:"input_preview"`
	OutputPreview string    `bson:"output_preview" json:"output_preview"`
	TokenUsage    int       `bson:"token_usage" json:"token_usage"`
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
	SessionID     string    `bson:"session_id" json:"session_id"`
	Status        string    `bson:"status" json:"status"`
}

type MetricsResponse struct {
	TotalRequests int64   `json:"total_requests"`
	AvgLatency    float64 `json:"avg_latency"`
	ErrorRate     float64 `json:"error_rate"`
}
