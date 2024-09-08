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

func (r *Repo) GetMessagesForFind(ctx context.Context, toSearch, chatID string) ([]int, error) {
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

func (r *Repo) GetMessageByMessageID(ctx context.Context, messageID int) (empty models.Message, err error) {
	query := fmt.Sprintf(`
	SELECT id, chat_id, message_id, text
	FROM %s
	WHERE message_id = $1
`, r.tableName)

	row := r.db.QueryRowContext(ctx, query, messageID)
	var msg models.Message

	if err = row.Scan(&msg.ID, &msg.ChatID, &msg.MessageID, &msg.Text); err != nil {
		return empty, fmt.Errorf("query get message: %w", err)
	}

	if err = row.Err(); err != nil {
		return empty, fmt.Errorf("invalid row: %w", err)
	}

	return msg, nil
}
