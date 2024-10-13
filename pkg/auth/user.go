package auth

import (
	"context"
	"rag-demo/types"
	"sync"
	"os"
	"github.com/google/uuid"
	"github.com/go-chi/jwtauth/v5"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	CreateUser(ctx context.Context, userName string, resultCh types.ResultChannel, wg *sync.WaitGroup)
	GetUser(ctx context.Context, userID uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup)
	GenerateJWT(ctx context.Context, session types.Session) (string, error)
	ValidateJWT(ctx context.Context, token string) (types.Session, error)
}

// AuthServiceImpl is the implementation of AuthService.
type AuthServiceImpl struct {
	UserGateway types.UserTableGateway
	SessionGateway types.SessionTableGateway
	tokenAuth *jwtauth.JWTAuth
}

// NewAuthService creates a new instance of AuthServiceImpl.
func NewAuthService(userGateway types.UserTableGateway, sessionGateway types.SessionTableGateway) AuthService {
	secret := os.Getenv("JWT_SECRET")
    algorithm := os.Getenv("JWT_ALGORITHM")

    // refactor to add better error handling 
    if secret == "" {
        return nil 
    }

    if algorithm == "" {
        return nil
    }
	tokenAuth := jwtauth.New(algorithm, []byte(secret), nil)
	return &AuthServiceImpl{UserGateway: userGateway, SessionGateway: sessionGateway, tokenAuth: tokenAuth}
}

// CreateUser creates a new user with the given name using the Go handler pattern.
// In auth/user.go
func (as *AuthServiceImpl) CreateUser(ctx context.Context, userName string, resultCh types.ResultChannel, wg *sync.WaitGroup) {
    defer wg.Done()

    newUser := types.User{
        UserID: uuid.New(),
        Name:   userName,
    }

    success, err := as.UserGateway.CreateUser(ctx, newUser)
    if err != nil || !success {
        resultCh <- types.Result{
            Data:    types.User{},
            Error:   err,
            Success: false,
        }
        return
    }

    // Create a new session
    newSession := types.Session{
        ID:     uuid.New(),
        UserID: newUser.UserID,
    }

    success, err = as.SessionGateway.CreateSession(ctx, newSession)
    if err != nil || !success {
        resultCh <- types.Result{
            Data:    types.User{},
            Error:   err,
            Success: false,
        }
        return
    }

    // Generate JWT
    token, err := as.GenerateJWT(ctx, newSession)
    if err != nil {
        resultCh <- types.Result{
            Data:    types.User{},
            Error:   err,
            Success: false,
        }
        return
    }

    resultCh <- types.Result{
        Data:    types.CreateUserResult{User: newUser, Token: token},
        Error:   nil,
        Success: true,
    }
}

func (as *AuthServiceImpl) GetUser(ctx context.Context, userID uuid.UUID, resultCh types.ResultChannel, wg *sync.WaitGroup) {
	defer wg.Done()

	user, err := as.UserGateway.GetUser(ctx, userID)
	if err != nil {
		resultCh <- types.Result{
			Data:    types.User{},
			Error:   err,
			Success: false,
		}
		return
	}

	resultCh <- types.Result{
		Data:    user,
		Error:   nil,
		Success: true,
	}
}