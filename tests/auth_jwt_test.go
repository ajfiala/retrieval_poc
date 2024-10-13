package tests

import (
    "context"
    "testing"
    "rag-demo/pkg/auth"
    "rag-demo/pkg/db"
    // "rag-demo/pkg/session"
    "rag-demo/types"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
    "github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
    dbURL := "postgresql://myuser:mypassword@localhost:5432/goragdb"
    testDBPool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    return testDBPool
}

func TestGenerateJWT(t *testing.T) {
    ctx := context.Background()
    dbPool := setupTestDB(t)
    defer dbPool.Close()

    godotenv.Load("../.env")

    sessionGateway := db.NewSessionTableGateway(dbPool)
    userGateway := db.NewUserTableGateway(dbPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    user := types.User{
        UserID: uuid.New(),
        Name:   "Test User",
    }

    session := types.Session{
        ID: uuid.New(),
        UserID:    user.UserID,
    }

    token, err := authService.GenerateJWT(ctx, session)

    assert.NoError(t, err, "GenerateJWT should not return an error")
    assert.NotEmpty(t, token, "Generated token should not be empty")

    // Optionally, you can add more assertions about the token's structure
}

func TestValidateJWT2(t *testing.T) {
    ctx := context.Background()
    dbPool := setupTestDB(t)
    defer dbPool.Close()

    // Load environment variables
    godotenv.Load("../.env")

    userGateway := db.NewUserTableGateway(dbPool)
    sessionGateway := db.NewSessionTableGateway(dbPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    user := types.User{
        UserID: uuid.New(),
        Name:   "Test User",
    }

    session := types.Session{
        ID: uuid.New(),
        UserID:    user.UserID,
    }


    // **Insert the user into the database**
    success, err := userGateway.CreateUser(ctx, user)
    assert.NoError(t, err, "CreateUser should not return an error")
    assert.True(t, success, "CreateUser should return true")

    // **Insert the session into the database**
    success, err = sessionGateway.CreateSession(ctx, session)
    assert.NoError(t, err, "CreateSession should not return an error")
    assert.True(t, success, "CreateSession should return true")


    // Ensure the user is deleted after the test
    defer func() {
        success, err := userGateway.DeleteUser(ctx, user.UserID)
        assert.NoError(t, err, "DeleteUser should not return an error")
        assert.True(t, success, "DeleteUser should return true")
    }()

    // Generate token for the user
    token, err := authService.GenerateJWT(ctx, session)
    assert.NoError(t, err, "GenerateJWT should not return an error")
	
	// Validate the token
	validatedSession, err := authService.ValidateJWT(ctx, token)
	assert.NoError(t, err, "ValidateJWT should not return an error")
	assert.Equal(t, user.UserID, validatedSession.UserID, "UserID should match")
	assert.Equal(t, session.ID, validatedSession.ID, "session ID should match")
}
