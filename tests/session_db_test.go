package tests

import (
	"context"
	"rag-demo/pkg/db"
	"rag-demo/types"

	// "rag-demo/types"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
    // Load environment variables from .env file if present
    if err := godotenv.Load("../.env"); err != nil {
        // Handle error if .env file is not found
        // For testing, we can set default environment variables
		fmt.Printf("Error loading .env file")
    }
}

func TestSessionTableGateway(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

	sessionGateway := db.NewSessionTableGateway(pool)

    userGateway := db.NewUserTableGateway(pool)

    testUser := types.User{
        UserID: uuid.New(),
        Name:   "Don Pizza",
    }

	testSession := types.Session{
		ID: uuid.New(),
		UserID: testUser.UserID,
	}

	success, err := userGateway.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if !success {
		t.Fatalf("CreateUser returned false")
	}
	
	t.Run("CreateSession", func(t *testing.T) {
        success, err := sessionGateway.CreateSession(ctx, testSession)
        if err != nil {
            t.Fatalf("CreateSession failed: %v", err)
        }
        if !success {
            t.Fatalf("CreateSession returned false")
        }
    })
	
	t.Run("GetSession", func(t *testing.T) {
		session, err := sessionGateway.GetSession(ctx, testSession.ID)
		if err != nil {
			t.Fatalf("GetSession failed: %v", err)
		}
		if session.ID != testSession.ID {
			t.Errorf("Expected ID %v, got %v", testSession.ID, session.ID)
		}
		if session.UserID != testSession.UserID {
			t.Errorf("Expected UserID %v, got %v", testSession.UserID, session.UserID)
		}
	})

}