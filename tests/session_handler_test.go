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
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

func TestCreateSessionHandler(t *testing.T) {
    router := chi.NewRouter()
    testDBPool, err := pgxpool.New(context.Background(), "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    userGateway := db.NewUserTableGateway(testDBPool)
    sessionGateway := db.NewSessionTableGateway(testDBPool)
    sessionService := message.NewSessionService(sessionGateway)

    router.Post("/api/v1/session", handlers.HandleCreateSession(sessionService))

    // Create a test user
    newUser := types.User{
        UserID: uuid.New(),
        Name:   "testuser",
    }

    _, err = userGateway.CreateUser(context.Background(), newUser)
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }

    newSessionRequest := types.NewSessionRequest{
        UserID: newUser.UserID,
    }

    body, err := json.Marshal(newSessionRequest)
    if err != nil {
        t.Fatalf("Failed to marshal new session request: %v", err)
    }

    req, err := http.NewRequest("POST", "/api/v1/session", bytes.NewReader(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var session types.Session
    err = json.NewDecoder(rr.Body).Decode(&session)
    if err != nil {
        t.Errorf("Failed to decode response body: %v", err)
    }

    // Verify the session's UserID
    if session.UserID != newUser.UserID {
        t.Errorf("Session has incorrect UserID: got %v, want %v", session.UserID, newUser.UserID)
    }

    _, err = userGateway.DeleteUser(context.Background(), newUser.UserID)
    if err != nil {
        t.Errorf("Failed to delete test user: %v", err)
    }
}
