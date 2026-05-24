package repository

import (
	"context"
	"time"

	"chatbot/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
	GetMessages(ctx context.Context, sessionID string) ([]models.Message, error)
	GetAllMessages(ctx context.Context, sessionID string) ([]models.Message, error)
	AppendMessage(ctx context.Context, sessionID string, msg models.Message) error
	ListSessions(ctx context.Context) ([]string, error)
}

type chatRepository struct {
	col *mongo.Collection
}

func NewChatRepository(db *mongo.Database) ChatRepository {
	return &chatRepository{col: db.Collection("chats")}
}

// GetMessages returns the last 5 messages — used for LLM context.
func (r *chatRepository) GetMessages(ctx context.Context, sessionID string) ([]models.Message, error) {
	msgs, err := r.GetAllMessages(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if len(msgs) > 5 {
		msgs = msgs[len(msgs)-5:]
	}
	return msgs, nil
}

// GetAllMessages returns the full conversation history — used for UI restore.
func (r *chatRepository) GetAllMessages(ctx context.Context, sessionID string) ([]models.Message, error) {
	var chat models.Chat
	err := r.col.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&chat)
	if err == mongo.ErrNoDocuments {
		return []models.Message{}, nil
	}
	if err != nil {
		return nil, err
	}
	return chat.Messages, nil
}

// ListSessions returns all session IDs sorted by most recently created.
func (r *chatRepository) ListSessions(ctx context.Context) ([]string, error) {
	opts := options.Find().
		SetProjection(bson.M{"session_id": 1}).
		SetSort(bson.M{"_id": -1}).
		SetLimit(20)

	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []string
	for cursor.Next(ctx) {
		var doc struct {
			SessionID string `bson:"session_id"`
		}
		if err := cursor.Decode(&doc); err == nil && doc.SessionID != "" {
			sessions = append(sessions, doc.SessionID)
		}
	}
	if sessions == nil {
		sessions = []string{}
	}
	return sessions, nil
}

func (r *chatRepository) AppendMessage(ctx context.Context, sessionID string, msg models.Message) error {
	msg.Timestamp = time.Now().UTC()
	filter := bson.M{"session_id": sessionID}
	update := bson.M{
		"$push":        bson.M{"messages": msg},
		"$setOnInsert": bson.M{"session_id": sessionID},
	}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}
