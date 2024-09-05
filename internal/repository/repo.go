package repository

import (
	"context"
	"database/sql"
	"fmt"

	"transcripter_bot/internal/models"
)

type (
	Repo struct {
		db        *sql.DB
		tableName string
	}
)

func New(db *sql.DB, tableName string) *Repo {
	return &Repo{
		db:        db,
		tableName: tableName,
	}
}

func (r *Repo) GetMessages(ctx context.Context, toSearch, chatID string) ([]int, error) {
	query := fmt.Sprintf(`
SELECT message_id
FROM %s
WHERE search_vector @@ to_tsquery('english', $1) AND  chat_id = $2
ORDER BY created_at DESC
;`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, toSearch, chatID)
	if err != nil {
		return nil, fmt.Errorf("query messages: %w", err)
	}
	defer rows.Close()

	var result []int
	for rows.Next() {

		var messageID int

		if err := rows.Scan(&messageID); err != nil {
			return nil, fmt.Errorf("scan messages: %w", err)
		}

		result = append(result, messageID)
	}

	return result, nil
}

func (r *Repo) SaveMessage(ctx context.Context, message models.Message) error {

	query := fmt.Sprintf(`
	INSERT INTO %s
	(chat_id, message_id, text, search_vector) 
	VALUES 
	($1, $2, $3, to_tsvector($4))
`, r.tableName)

	_, err := r.db.ExecContext(ctx, query, message.ChatID, message.MessageID, message.Text, message.Text)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}

	return nil
}

//func (s Repo) GetMessages(ctx context.Context, target string, chatID int64) ([]int64, error) {
//	filter := bson.M{
//		"text": bson.M{
//			"$regex":   target,
//			"$options": "i",
//		},
//		"chat_id": chatID,
//	}
//
//	projection := bson.M{
//		"message_id": 1,
//	}
//
//	cursor, err := s.collection.Find(ctx, filter, options.Find().SetProjection(projection))
//	if err != nil {
//		return nil, fmt.Errorf("failed to find in collection: %w", err)
//	}
//	defer cursor.Close(ctx)
//
//	var results []getMessagesResponse
//
//	if err = cursor.All(ctx, &results); err != nil {
//		return nil, fmt.Errorf("failed to decode messages: %w", err)
//	}
//
//	matchingIDs := make([]int64, 0, len(results))
//	for _, result := range results {
//		matchingIDs = append(matchingIDs, result.MessageID)
//	}
//
//	return matchingIDs, nil
//}
//
//func (s *Repo) SaveMessage(ctx context.Context, message models.Message) error {
//	if _, err := s.collection.InsertOne(ctx, message); err != nil {
//		return fmt.Errorf("failed to save message: %w", err)
//	}
//
//	return nil
//}
