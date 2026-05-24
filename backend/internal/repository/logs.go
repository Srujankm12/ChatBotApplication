package repository

import (
	"context"

	"chatbot/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository interface {
	InsertLog(ctx context.Context, log models.InferenceLog) error
	GetMetrics(ctx context.Context) (models.MetricsResponse, error)
}

type logRepository struct {
	col *mongo.Collection
}

func NewLogRepository(db *mongo.Database) LogRepository {
	return &logRepository{col: db.Collection("logs")}
}

func (r *logRepository) InsertLog(ctx context.Context, log models.InferenceLog) error {
	_, err := r.col.InsertOne(ctx, log)
	return err
}

func (r *logRepository) GetMetrics(ctx context.Context) (models.MetricsResponse, error) {
	total, err := r.col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return models.MetricsResponse{}, err
	}

	errorCount, err := r.col.CountDocuments(ctx, bson.M{"status": "error"})
	if err != nil {
		return models.MetricsResponse{}, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":        nil,
			"avg_latency": bson.M{"$avg": "$latency_ms"},
		}}},
	}
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return models.MetricsResponse{}, err
	}
	defer cursor.Close(ctx)

	var avgLatency float64
	if cursor.Next(ctx) {
		var result struct {
			AvgLatency float64 `bson:"avg_latency"`
		}
		if err := cursor.Decode(&result); err == nil {
			avgLatency = result.AvgLatency
		}
	}

	var errorRate float64
	if total > 0 {
		errorRate = float64(errorCount) / float64(total)
	}

	return models.MetricsResponse{
		TotalRequests: total,
		AvgLatency:    avgLatency,
		ErrorRate:     errorRate,
	}, nil
}
