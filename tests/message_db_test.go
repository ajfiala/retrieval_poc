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

func TestMessageTableGateway(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

	sessionGateway := db.NewSessionTableGateway(pool)

    userGateway := db.NewUserTableGateway(pool)

	messageGateway := db.NewMessageTableGateway(pool)

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
	
	t.Run("Create Message", func(t *testing.T) {
		testMessage := types.Message{
			ID: uuid.New(),
			UserId: testUser.UserID,
			SessionId: testSession.ID,
			UserMessage: "Hello",
			AiMessage: "Hi",
		}

		success, err := messageGateway.StoreMessage(ctx, testMessage, testSession)
		if err != nil {
			t.Fatalf("StoreMessage failed: %v", err)
		}
		if !success {
			t.Fatalf("StoreMessage returned false")
		}
	})			


	t.Run("List Messages", func(t *testing.T) {
		MessagesList, err := messageGateway.RetrieveMessages(ctx, testSession.ID)
		fmt.Println(MessagesList)
		if err != nil {
			t.Fatalf("StoreMessage failed: %v", err)
		}
		if len(MessagesList.Messages) == 0 {
			t.Fatalf("RetrieveMessages returned nothing")
		}
		if len(MessagesList.Messages) != 1 {
			t.Fatalf("RetrieveMessages returned incorrect number of messages")
		}
	})			
}