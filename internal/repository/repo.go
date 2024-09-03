package repository

import (
	"context"
	"fmt"
	"log/slog"

	"transcripter_bot/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func New(db *mongo.Database, collectionName string, log *slog.Logger) *Repo {
	collection := db.Collection(collectionName)

	return &Repo{
		collection: collection,
		log:        log,
	}
}

type getMessagesResponse struct {
	MessageID int64 `bson:"message_id"`
}

func (s Repo) GetMessages(ctx context.Context, target string, chatID int64) ([]int64, error) {
	filter := bson.M{
		"text": bson.M{
			"$regex":   target,
			"$options": "i",
		},
		"chat_id": chatID,
	}

	projection := bson.M{
		"message_id": 1,
	}

	cursor, err := s.collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to find in collection: %w", err)
	}
	defer cursor.Close(ctx)

	var results []getMessagesResponse

	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	matchingIDs := make([]int64, 0, len(results))
	for _, result := range results {
		matchingIDs = append(matchingIDs, result.MessageID)
	}

	return matchingIDs, nil
}

func (s *Repo) SaveMessage(ctx context.Context, message models.Message) error {
	if _, err := s.collection.InsertOne(ctx, message); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}
