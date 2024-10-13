package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rag-demo/pkg/db"
	"rag-demo/pkg/handlers"
	"rag-demo/pkg/auth"
	"github.com/google/uuid"
	"rag-demo/pkg/message"
	"rag-demo/types"
	"testing"
	"fmt"
	"time"
	"github.com/go-chi/chi/v5"
	// "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func init() {
    // Load environment variables from .env file if present
    if err := godotenv.Load("../.env"); err != nil {
        // Handle error if .env file is not found
        // For testing, we can set default environment variables
		fmt.Printf("Error loading .env file")
    }
}

func TestSendMessageHandler(t *testing.T) {
	router := chi.NewRouter()
	ctx := context.Background()
	
	// Set a timeout for the entire test
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	testDBPool, err := pgxpool.New(ctx, "postgresql://myuser:mypassword@localhost:5432/goragdb")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer testDBPool.Close()

	userGateway := db.NewUserTableGateway(testDBPool)
	sessionGateway := db.NewSessionTableGateway(testDBPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    user := types.User{
        UserID: uuid.New(),
        Name:   "Test User",
    }

	// store user 
	ok, err := userGateway.CreateUser(ctx, user)
	if !ok {
		t.Fatalf("Failed to create user: %v", err)
	}

    session := types.Session{
        ID: uuid.New(),
        UserID:    user.UserID,
    }

	// store session
	ok, err = sessionGateway.CreateSession(ctx, session)
	if !ok {
		t.Fatalf("Failed to create session: %v", err)
	}

	// create a JWT token that we will use as cookie in request
	token, err := authService.GenerateJWT(ctx, session)
	if err != nil {
		t.Fatalf("Failed to generate JWT token: %v", err)
	}

	// validate the JWT 
	validatedSession, err := authService.ValidateJWT(ctx, token)
	if err != nil {
		t.Fatalf("Failed to validate JWT token: %v", err)
	}

	if validatedSession.UserID != user.UserID {
		t.Fatalf("Failed to validate JWT token: %v", err)
	}
	if validatedSession.ID != session.ID {
		t.Fatalf("Failed to validate JWT token: %v", err)
	}
	



	messageGateway := db.NewMessageTableGateway(testDBPool)

	// Create the real Bedrock service
	bedrockService, err := message.NewBedrockRuntimeService("anthropic")
	if err != nil {
		t.Fatalf("Failed to create Bedrock service: %v", err)
	}
	messageService := message.NewMessageService(messageGateway, bedrockService)

	router.Post("/api/v1/message", handlers.HandleSendMessage(authService.(*auth.AuthServiceImpl), messageService))


	// // Create a test message request
	testMessage := types.MessageRequest{
		Text:      "Hello, AI!",
	}

	body, err := json.Marshal(testMessage)
	if err != nil {
		t.Fatalf("Failed to marshal message request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/v1/message", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add the JWT token as a cookie
	cookie := &http.Cookie{
		Name:  "access-token",
		Value: token, // Add "Bearer " prefix
		Path:  "/",
	}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var response types.Message
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "Failed to decode response body")

	fmt.Printf("Response: %v\n", response)
	
	assert.Equal(t, testMessage.Text, response.UserMessage, "Response has incorrect UserMessage")
	assert.NotEmpty(t, response.AiMessage, "AI Message should not be empty")

	// Clean up the test message
	// _, err = messageGateway.DeleteMessage(ctx, response.ID)
	// assert.NoError(t, err, "Failed to delete test message")
}