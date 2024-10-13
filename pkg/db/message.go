package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"rag-demo/types"
	"fmt"
	"time"
)

// UserTableGatewayImpl is the implementation of UserTableGateway using pgxpool.
type MessageTableGatewayImpl struct {
	Pool *pgxpool.Pool
}

// NewUserTableGateway creates a new instance of UserTableGatewayImpl.
func NewMessageTableGateway(pool *pgxpool.Pool) types.MessageTableGateway {
	return &MessageTableGatewayImpl{Pool: pool}
}

func (mtg *MessageTableGatewayImpl) StoreMessage(ctx context.Context, message types.Message, session types.Session) (bool, error) {
    fmt.Println("Entering StoreMessage function")
    
    // Create a new context with a timeout
    ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    fmt.Println("Executing database query")
    _, err := mtg.Pool.Exec(ctxWithTimeout, "INSERT INTO message (uuid, user_id, session_id, user_message, ai_message) VALUES ($1, $2, $3, $4, $5)",
        message.ID,
        session.UserID,
        session.ID,
        message.UserMessage,
        message.AiMessage)
    
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return false, err
    }
    
    fmt.Println("StoreMessage completed successfully")
    return true, nil
}

func (mtg *MessageTableGatewayImpl) RetrieveMessages(ctx context.Context, session_id uuid.UUID) (types.MessageList, error) {
	rows, err := mtg.Pool.Query(ctx, "SELECT uuid, user_id, session_id, user_message, ai_message FROM message WHERE session_id = $1", session_id)
	if err != nil {
		return types.MessageList{}, err
	}
	defer rows.Close()

	var messages []types.Message
	for rows.Next() {
		var message types.Message
		err := rows.Scan(&message.ID, &message.UserId, &message.SessionId, &message.UserMessage, &message.AiMessage)
		if err != nil {
			return types.MessageList{}, err
		}
		messages = append(messages, message)
	}
	return types.MessageList{Messages: messages}, nil
}