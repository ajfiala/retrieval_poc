package tests

import (
    "context"
    "testing"

    "rag-demo/pkg/db"
    "rag-demo/pkg/message"
    "rag-demo/types"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
    "github.com/stretchr/testify/assert"
)

// setupTestDB initializes a test database connection.
func setupTestSessionDB(t *testing.T) *pgxpool.Pool {
    dbURL := "postgresql://myuser:mypassword@localhost:5432/goragdb"
    testDBPool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    return testDBPool
}


func TestSessionService_CreateSession(t *testing.T) {
    ctx := context.Background()
    dbPool := setupTestSessionDB(t)
    defer dbPool.Close()

    // Load environment variables
    godotenv.Load("../.env")

    sessionGateway := db.NewSessionTableGateway(dbPool)
    sessionService := message.NewSessionService(sessionGateway)

    userGateway := db.NewUserTableGateway(dbPool)

    // Test data
    testUser := types.User{
        UserID: uuid.New(),
        Name:   "Don Pizza",
    }

    // Create User to prepare for session test
	success, err := userGateway.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if !success {
		t.Fatalf("CreateUser returned false")
	}

    // Create a new session
    session, err := sessionService.CreateSession(ctx, testUser.UserID)

    assert.NoError(t, err, "CreateSession should not return an error")
    assert.NotNil(t, session, "Session should not be nil")
    assert.Equal(t, testUser.UserID, session.UserID, "UserID should match")
    assert.NotEqual(t, uuid.Nil, session.ID, "Session ID should not be nil")

}

func TestSessionService_GetSession(t *testing.T) {
	ctx := context.Background()
	dbPool := setupTestSessionDB(t)
	defer dbPool.Close()

	// Load environment variables
	godotenv.Load("../.env")

	sessionGateway := db.NewSessionTableGateway(dbPool)
	sessionService := message.NewSessionService(sessionGateway)

	userGateway := db.NewUserTableGateway(dbPool)

	// Test data
	testUser := types.User{
		UserID: uuid.New(),
		Name:   "Don Pizza",
	}

	// Create User to prepare for session test
	success, err := userGateway.CreateUser(ctx, testUser)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if !success {
		t.Fatalf("CreateUser returned false")
	}

	// Create a new session
	session, err := sessionService.CreateSession(ctx, testUser.UserID)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Get the session
	retrievedSession, err := sessionService.GetSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	assert.Equal(t, session.ID, retrievedSession.ID, "Session ID should match")
	assert.Equal(t, session.UserID, retrievedSession.UserID, "UserID should match")
}