package types

import (
	"context"
	"github.com/google/uuid"
)

// represents the session object in the database
type Session struct {
	ID uuid.UUID  `json:"session_id"`
	UserID uuid.UUID `json:"user_id"`
}

// represents the payload for starting a new session for a user
type NewSessionRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type SessionTableGateway interface {
	CreateSession(ctx context.Context, session Session) (bool, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (Session, error)
}