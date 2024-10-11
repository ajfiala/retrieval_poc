package tests

import (
    "context"
    "rag-demo/pkg/db"
    "rag-demo/types"
    "testing"
	"fmt"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
    "os"
)

func init() {
    // Load environment variables from .env file if present
    if err := godotenv.Load("../.env"); err != nil {

		fmt.Printf("Error loading .env file")
    }
}

func TestUserTableGateway(t *testing.T) {
    ctx := context.Background()
    // Initialize the database connection
    pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

    userGateway := db.NewUserTableGateway(pool)

    testUser := types.User{
        UserID: uuid.New(),
        Name:   "Test User",
    }

    t.Run("CreateUser", func(t *testing.T) {
        success, err := userGateway.CreateUser(ctx, testUser)
        if err != nil {
            t.Fatalf("CreateUser failed: %v", err)
        }
        if !success {
            t.Fatalf("CreateUser returned false")
        }
    })

    t.Run("GetUser", func(t *testing.T) {
        user, err := userGateway.GetUser(ctx, testUser.UserID)
        if err != nil {
            t.Fatalf("GetUser failed: %v", err)
        }
        if user.UserID != testUser.UserID {
            t.Errorf("Expected UserID %v, got %v", testUser.UserID, user.UserID)
        }
        if user.Name != testUser.Name {
            t.Errorf("Expected Name %v, got %v", testUser.Name, user.Name)
        }
    })

    t.Run("DeleteUser", func(t *testing.T) {
        success, err := userGateway.DeleteUser(ctx, testUser.UserID)
        if err != nil {
            t.Fatalf("DeleteUser failed: %v", err)
        }
        if !success {
            t.Fatalf("DeleteUser returned false")
        }

        _, err = userGateway.GetUser(ctx, testUser.UserID)
        if err == nil {
            t.Fatalf("Expected error when fetching deleted user, got nil")
        }
    })
}
