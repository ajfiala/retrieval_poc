// types/user.go
package types

import (
	"context"
	"github.com/google/uuid"
	// "github.com/go-playground/validator/v10"
)

// User represents a user in the system.
type User struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

// NewUserRequest represents the payload for creating a new user.
type NewUserRequest struct {
	Name string `json:"name" validate:"required"`
}

// UserTableGateway defines the interface for user-related database operations.
type UserTableGateway interface {
	CreateUser(ctx context.Context, user User) (bool, error)
	GetUser(ctx context.Context, userID uuid.UUID) (User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) (bool, error)
}
