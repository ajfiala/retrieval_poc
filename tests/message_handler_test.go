package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rag-demo/pkg/db"
	"rag-demo/pkg/handlers"
	"rag-demo/pkg/message"
	"rag-demo/types"
	"testing"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func sageHandler(t *testing.T) {
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

	messageGateway := db.NewMessageTableGateway(testDBPool)

	// Create the real Bedrock service
	bedrockService, err := message.NewBedrockRuntimeService("anthropic")
	if err != nil {
		t.Fatalf("Failed to create Bedrock service: %v", err)
	}
	messageService := message.NewMessageService(messageGateway, bedrockService)

	router.Post("/api/v1/message", handlers.HandleSendMessage(messageService))

	// Create a test message request
	testMessage := types.MessageRequest{
		Text:      "Hello, AI!",
		SessionId: uuid.New(),
		UserId:    uuid.New(),
	}

	body, err := json.Marshal(testMessage)
	if err != nil {
		t.Fatalf("Failed to marshal message request: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/message", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var response types.Message
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "Failed to decode response body")

	assert.Equal(t, testMessage.UserId, response.UserId, "Response has incorrect UserID")
	assert.Equal(t, testMessage.SessionId, response.SessionId, "Response has incorrect SessionID")
	assert.Equal(t, testMessage.Text, response.UserMessage, "Response has incorrect UserMessage")
	assert.NotEmpty(t, response.AiMessage, "AI Message should not be empty")

	// Clean up the test message
	_, err = messageGateway.DeleteMessage(ctx, response.ID)
	assert.NoError(t, err, "Failed to delete test message")
}