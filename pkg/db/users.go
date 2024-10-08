package db

import (
	"context"
	"rag-demo/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserTableGatewayImpl is the implementation of UserTableGateway using pgxpool.
type UserTableGatewayImpl struct {
	Pool *pgxpool.Pool
}

// NewUserTableGateway creates a new instance of UserTableGatewayImpl.
func NewUserTableGateway(pool *pgxpool.Pool) types.UserTableGateway {
	return &UserTableGatewayImpl{Pool: pool}
}

// CreateUser inserts a new user into the database.
func (utg *UserTableGatewayImpl) CreateUser(ctx context.Context, user types.User) (bool, error) {
	_, err := utg.Pool.Exec(ctx, "INSERT INTO users (uuid, name, active) VALUES ($1, $2, $3)", user.UserID, user.Name, true)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUser retrieves a user by userID.
func (utg *UserTableGatewayImpl) GetUser(ctx context.Context, userID uuid.UUID) (types.User, error) {
	var user types.User
	// userIdStr := userID.String()
	err := utg.Pool.QueryRow(ctx, "SELECT uuid, name FROM users WHERE uuid = $1", userID).Scan(&user.UserID, &user.Name)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}

// DeleteUser removes a user by userID.
func (utg *UserTableGatewayImpl) DeleteUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	commandTag, err := utg.Pool.Exec(ctx, "DELETE FROM users WHERE uuid = $1", userID)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() == 0 {
		return false, nil // No rows deleted
	}
	return true, nil
}
