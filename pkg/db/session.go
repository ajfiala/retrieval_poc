package db

import (
	"context"
	"rag-demo/types"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

// UserTableGatewayImpl is the implementation of UserTableGateway using pgxpool.
type SessionTableGatewayImpl struct {
	Pool *pgxpool.Pool
}

// NewUserTableGateway creates a new instance of UserTableGatewayImpl.
func NewSessionTableGateway(pool *pgxpool.Pool) types.SessionTableGateway {
	return &SessionTableGatewayImpl{Pool: pool}
}

func (stg *SessionTableGatewayImpl) CreateSession(ctx context.Context, session types.Session) (bool, error) {
	fmt.Println("Creating session: ", session)
	_, err := stg.Pool.Exec(ctx, "INSERT INTO session (uuid, user_id, active) VALUES ($1, $2, $3)", session.ID, session.UserID, true)
	if err != nil {
		return false, err
	}
	fmt.Println("Session created: ", session)
	return true, nil
}

func (stg *SessionTableGatewayImpl) GetSession(ctx context.Context, sessionID uuid.UUID) (types.Session, error) {
	var session types.Session
	err := stg.Pool.QueryRow(ctx, "SELECT uuid, user_id FROM session WHERE uuid = $1", sessionID).Scan(&session.ID, &session.UserID)
	if err != nil {
		return types.Session{}, err
	}
	return session, nil
}