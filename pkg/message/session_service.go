package message 

import (
	"context"
	"rag-demo/types"
	// google uuid
	"sync"
	"github.com/google/uuid"
)

// SessionService defines the interface for session-related operations.
type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup) 
	GetSession(ctx context.Context, sessionID uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup) 
}

type SessionServiceImpl struct {
	SessionGateway types.SessionTableGateway
}

func NewSessionService(sessionGateway types.SessionTableGateway) SessionService {
	return &SessionServiceImpl{SessionGateway: sessionGateway}
}

func (ss *SessionServiceImpl) CreateSession(ctx context.Context, userID uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup) {
    defer wg.Done()

    newSession := types.Session{
        ID:     uuid.New(),
        UserID: userID,
    }

    success, err := ss.SessionGateway.CreateSession(ctx, newSession)
    if err != nil || !success {
        resultCh <- types.Result{
            Data:    nil,
            Error:   err,
            Success: false,
        }
        return
    }

    resultCh <- types.Result{
        Data:    newSession,
        Error:   nil,
        Success: true,
    }
}

func (ss *SessionServiceImpl) GetSession(ctx context.Context, sessionID uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup) {
    defer wg.Done()

    session, err := ss.SessionGateway.GetSession(ctx, sessionID)
    if err != nil {
        resultCh <- types.Result{
            Data:    nil,
            Error:   err,
            Success: false,
        }
        return
    }

    resultCh <- types.Result{
        Data:    session,
        Error:   nil,
        Success: true,
    }
}