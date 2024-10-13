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
	"rag-demo/pkg/assistant"
	"rag-demo/types"
	"testing"
	"fmt"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("Error loading .env file")
	}
}

func TestHandleCreateAssistant(t *testing.T) {
	router := chi.NewRouter()
	ctx := context.Background()
	
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

	assistantGateway := db.NewAssistantTableGateway(testDBPool)
	assistantService := assistant.NewAssistantService(assistantGateway)

	user := types.User{
		UserID: uuid.New(),
		Name:   "Test User",
	}

	ok, err := userGateway.CreateUser(ctx, user)
	if !ok {
		t.Fatalf("Failed to create user: %v", err)
	}

	session := types.Session{
		ID:     uuid.New(),
		UserID: user.UserID,
	}

	ok, err = sessionGateway.CreateSession(ctx, session)
	if !ok {
		t.Fatalf("Failed to create session: %v", err)
	}

	token, err := authService.GenerateJWT(ctx, session)
	if err != nil {
		t.Fatalf("Failed to generate JWT token: %v", err)
	}

	router.Post("/api/v1/assistant", handlers.HandleCreateAssistant(authService.(*auth.AuthServiceImpl), assistantService))

	testAssistant := types.Assistant{
		ID:       uuid.New(),
		Name:     "insurance assistant3",
		Type:     "rag",
		Config: types.Config{
			Provider: "anthropic",
			ModelId: "anthropic.claude-3-5-sonnet-20240620-v1:0",
		},
		SystemPrompts: 
`You are an insurance manager for a large Thai insurance company.
You are an expert on insurance products, rules, services, and policies. Users will
ask you questions about insurance products, rules, services, and policies. Respond
using the provided context. If the question is general enough, you can provide a
general answer. Otherwise, don't answer if you aren't sure or if the answer cannot be found
in the provided context. In this case, direct the user to another source of information.`,
		Metadata: types.Metadata{
			Title:       "Insurance Assistant",
			Description: "Assists users with insurance-related queries.",
			Icon:        "insurance_icon.png",
			Prompts:     []string{"How can I assist you with insurance today?",
		"How do I file a claim?",
		"What is the process for renewing my policy?",},
		},
	}

	body, err := json.Marshal(testAssistant)
	if err != nil {
		t.Fatalf("Failed to marshal assistant request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/v1/assistant", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	cookie := &http.Cookie{
		Name:  "access-token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var response types.Assistant
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "Failed to decode response body")

	assert.Equal(t, testAssistant.Name, response.Name, "Response has incorrect Name")
	assert.Equal(t, testAssistant.Type, response.Type, "Response has incorrect Type")
	assert.NotEmpty(t, response.ID, "Assistant ID should not be empty")

	// clean up with exec 
	_, err = testDBPool.Exec(ctx, "DELETE FROM assistant WHERE uuid = $1", response.ID)
	if err != nil {
		t.Fatalf("Failed to delete assistant: %v", err)
	}
	 
}

func TestHandleListAssistants(t *testing.T) {
	router := chi.NewRouter()
	ctx := context.Background()
	
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

	assistantGateway := db.NewAssistantTableGateway(testDBPool)
	assistantService := assistant.NewAssistantService(assistantGateway)

	user := types.User{
		UserID: uuid.New(),
		Name:   "Test User",
	}

	ok, err := userGateway.CreateUser(ctx, user)
	if !ok {
		t.Fatalf("Failed to create user: %v", err)
	}

	session := types.Session{
		ID:     uuid.New(),
		UserID: user.UserID,
	}

	testAssistant := types.Assistant{
		ID:       uuid.New(),
		Name:     "insurance assistant3",
		Type:     "rag",
		Config: types.Config{
			Provider: "anthropic",
			ModelId: "anthropic.claude-3-5-sonnet-20240620-v1:0",
		},
		SystemPrompts: 
`You are an insurance manager for a large Thai insurance company.
You are an expert on insurance products, rules, services, and policies. Users will
ask you questions about insurance products, rules, services, and policies. Respond
using the provided context. If the question is general enough, you can provide a
general answer. Otherwise, don't answer if you aren't sure or if the answer cannot be found
in the provided context. In this case, direct the user to another source of information.`,
		Metadata: types.Metadata{
			Title:       "Insurance Assistant",
			Description: "Assists users with insurance-related queries.",
			Icon:        "insurance_icon.png",
			Prompts:     []string{"How can I assist you with insurance today?",
		"How do I file a claim?",
		"What is the process for renewing my policy?",},
		},
	}

	testAssistant2 := types.Assistant{
		ID:       uuid.New(),
		Name:     "insurance assistant4",
		Type:     "rag",
		Config: types.Config{
			Provider: "anthropic",
			ModelId: "anthropic.claude-3-5-sonnet-20240620-v1:0",
		},
		SystemPrompts: 
`You are an insurance manager for a large Thai insurance company.
You are an expert on insurance products, rules, services, and policies. Users will
ask you questions about insurance products, rules, services, and policies. Respond
using the provided context. If the question is general enough, you can provide a
general answer. Otherwise, don't answer if you aren't sure or if the answer cannot be found
in the provided context. In this case, direct the user to another source of information.`,
		Metadata: types.Metadata{
			Title:       "Insurance Assistant",
			Description: "Assists users with insurance-related queries.",
			Icon:        "insurance_icon.png",
			Prompts:     []string{"How can I assist you with insurance today?",
		"How do I file a claim?",
		"What is the process for renewing my policy?",},
		},
	}

	_, err = assistantGateway.CreateAssistant(ctx, testAssistant)
	if err != nil {
		t.Fatalf("CreateAssistant failed: %v", err)
	}

	_, err = assistantGateway.CreateAssistant(ctx, testAssistant2)
	if err != nil {
		t.Fatalf("CreateAssistant failed: %v", err)
	}

	ok, err = sessionGateway.CreateSession(ctx, session)
	if !ok {
		t.Fatalf("Failed to create session: %v", err)
	}

	token, err := authService.GenerateJWT(ctx, session)
	if err != nil {
		t.Fatalf("Failed to generate JWT token: %v", err)
	}

	router.Get("/api/v1/assistants", handlers.HandleListAssistants(authService.(*auth.AuthServiceImpl), assistantService))

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/assistants", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:  "access-token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var response types.AssistantList
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err, "Failed to decode response body")

	// Assert that the response is a list (even if empty)
	assert.NotNil(t, response.Assistants, "Assistants list should not be nil")

	// assert that length of AssistantList is 2 
	assert.Equal(t, 2, len(response.Assistants), "Assistants list should have 2 items")

	// clean up with exec 
	_, err = testDBPool.Exec(ctx, "DELETE FROM assistant WHERE uuid = $1", testAssistant.ID)
	if err != nil {
		t.Fatalf("Failed to delete assistant: %v", err)
	}
	_, err = testDBPool.Exec(ctx, "DELETE FROM assistant WHERE uuid = $1", testAssistant2.ID)
	if err != nil {
		t.Fatalf("Failed to delete assistant: %v", err)
	}
}