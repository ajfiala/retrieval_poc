package tests

import (
    "context"
    "testing"
	"sync"
    "rag-demo/pkg/db"
    "rag-demo/pkg/message"
    "rag-demo/types"
	"fmt"
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
	resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
	wg := &sync.WaitGroup{}
	wg.Add(1)

    sessionService.CreateSession(ctx, testUser.UserID, resultCh, wg)

	wg.Wait()             // Wait for the goroutine to finish
	result := <-resultCh  // Read the result from the channel

	fmt.Println("Result: ", result)
	assert.True(t, result.Success, "Result should be successful")

    assert.NoError(t, err, "CreateSession should not return an error")
	c := types.Session{
		ID: result.Data.(types.Session).ID,
		UserID: result.Data.(types.Session).UserID,
	}
    assert.NotNil(t, c, "Session should not be nil")
    assert.Equal(t, testUser.UserID, c.UserID, "UserID should match")
    assert.NotEqual(t, uuid.Nil, c.ID, "Session ID should not be nil")

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
	resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
	wg := &sync.WaitGroup{}
	wg.Add(1)

	sessionService.CreateSession(ctx, testUser.UserID, resultCh, wg)

	wg.Wait()             // Wait for the goroutine to finish
	result := <-resultCh  // Read the result from the channel

	assert.True(t, result.Success, "Result should be successful")
	newSession := types.Session{
		ID: result.Data.(types.Session).ID,
		UserID: result.Data.(types.Session).UserID,
	}
	
	resultCh = make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
	wg = &sync.WaitGroup{}
	wg.Add(1)
	// Get the session
	sessionService.GetSession(ctx, newSession.ID, resultCh, wg)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	wg.Wait()             // Wait for the goroutine to finish

	result = <-resultCh  // Read the result from the channel
	
	assert.True(t, result.Success, "Result should be successful")

	retrievedSession := types.Session{
		ID: result.Data.(types.Session).ID,
		UserID: result.Data.(types.Session).UserID,
	}

	assert.Equal(t, newSession.ID, retrievedSession.ID, "Session ID should match")
	assert.Equal(t, newSession.UserID, retrievedSession.UserID, "UserID should match")
}