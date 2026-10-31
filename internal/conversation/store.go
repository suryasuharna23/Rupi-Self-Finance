package conversation

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// IdempotencyStore records processed messages in processed_messages.
type IdempotencyStore struct{ db *pgxpool.Pool }

func NewIdempotencyStore(db *pgxpool.Pool) *IdempotencyStore { return &IdempotencyStore{db: db} }

// MarkProcessed returns true if the message has not been processed before.
func (s *IdempotencyStore) MarkProcessed(ctx context.Context, channel, messageID string) (bool, error) {
	tag, err := s.db.Exec(ctx, `
		INSERT INTO conversation_svc.processed_messages (channel, external_message_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING`, channel, messageID)
	return tag.RowsAffected() == 1, err
}

// Unmark reopens a message that failed so a retry can succeed.
func (s *IdempotencyStore) Unmark(ctx context.Context, channel, messageID string) error {
	_, err := s.db.Exec(ctx, `
		DELETE FROM conversation_svc.processed_messages
		WHERE channel = $1 AND external_message_id = $2`, channel, messageID)
	return err
}
