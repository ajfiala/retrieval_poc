package tests

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "rag-demo/pkg/api"
    "rag-demo/pkg/auth"
    "rag-demo/pkg/db"
    "rag-demo/pkg/handlers"
    "rag-demo/types"
    "sync"
    "testing"
	"fmt"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
)

// Initialize environment variables (e.g., JWT_SECRET, JWT_ALGORITHM)
func init() {
    if err := godotenv.Load("../.env"); err != nil {
        // Handle error if .env file is not found
		fmt.Printf("Error loading .env file")
    }
}

func TestHelloWorld(t *testing.T) {
    router := chi.NewRouter()
    router.Get("/", api.HelloWorld)

    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: Got %v, want %v", status, http.StatusOK)
    }

    expected := "Hello, World!"
    if rr.Body.String() != expected {
        t.Errorf("Handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
    }
}

func TestCreateUser(t *testing.T) {
    ctx := context.Background()
    testDBPool, err := pgxpool.New(ctx, "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    var wg sync.WaitGroup
    resultCh := make(types.ResultChannel)

    userGateway := db.NewUserTableGateway(testDBPool)
    sessionGateway := db.NewSessionTableGateway(testDBPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    wg.Add(1)
    go authService.CreateUser(ctx, "testuser", resultCh, &wg)

    result := <-resultCh
    wg.Wait()

    if !result.Success {
        t.Errorf("CreateUser operation failed: %v", result.Error)
    }

    user, ok := result.Data.(types.CreateUserResult)
    if !ok {
        t.Errorf("Result data is not of type User")
    }

    if user.User.UserID == uuid.Nil {
        t.Errorf("CreateUser returned User with invalid UserID: got %v", user.User.UserID)
    }

    // Clean up: delete the created user
    _, err = userGateway.DeleteUser(ctx, user.User.UserID)
    if err != nil {
        t.Errorf("Failed to delete test user: %v", err)
    }
}

func TestGetUser(t *testing.T) {
    ctx := context.Background()
    testDBPool, err := pgxpool.New(ctx, "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    userGateway := db.NewUserTableGateway(testDBPool)
    sessionGateway := db.NewSessionTableGateway(testDBPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    newUser := types.User{
        UserID: uuid.New(),
        Name:   "testuser",
    }

    _, err = userGateway.CreateUser(ctx, newUser)
    if err != nil {
        t.Errorf("Failed to create test user: %v", err)
    }

    var wg sync.WaitGroup
    resultCh := make(types.ResultChannel)

    wg.Add(1)
    go authService.GetUser(ctx, newUser.UserID, resultCh, &wg)

    result := <-resultCh
    wg.Wait()

    if !result.Success {
        t.Errorf("GetUser operation failed: %v", result.Error)
    }

    user, ok := result.Data.(types.User)
    if !ok {
        t.Errorf("Result data is not of type User")
    }

    if user.UserID != newUser.UserID {
        t.Errorf("GetUser returned incorrect user: got %v, want %v", user.UserID, newUser.UserID)
    }

    // Clean up: delete the created user
    _, err = userGateway.DeleteUser(ctx, newUser.UserID)
    if err != nil {
        t.Errorf("Failed to delete test user: %v", err)
    }
}

func TestCreateUserHandler(t *testing.T) {
    router := chi.NewRouter()
    testDBPool, err := pgxpool.New(context.Background(), "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    userGateway := db.NewUserTableGateway(testDBPool)
    sessionGateway := db.NewSessionTableGateway(testDBPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    router.Post("/api/v1/user", handlers.HandleCreateUser(authService))

    newUser := types.NewUserRequest{Name: "testuser"}
    body, _ := json.Marshal(newUser)
    req, err := http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var createdUser types.User
    err = json.NewDecoder(rr.Body).Decode(&createdUser)
    if err != nil {
        t.Errorf("Failed to decode response body: %v", err)
    }

    if createdUser.UserID == uuid.Nil {
        t.Errorf("Created user has invalid UserID")
    }

    if createdUser.Name != "testuser" {
        t.Errorf("Created user has incorrect name: got %v, want %v", createdUser.Name, "testuser")
    }

    // Check if the JWT token is set in the header
    tokenHeader := rr.Header().Get("access-token")
    if tokenHeader == "" {
        t.Errorf("access-token header is not set")
    }

    // Check if the JWT token is set as a cookie
    cookies := rr.Result().Cookies()
    var tokenCookie *http.Cookie
    for _, cookie := range cookies {
        if cookie.Name == "access-token" {
            tokenCookie = cookie
            break
        }
    }
    if tokenCookie == nil {
        t.Errorf("access-token cookie is not set")
    } else {
        if tokenCookie.Value == "" {
            t.Errorf("access-token cookie has no value")
        }
        // Optionally, validate the JWT token here
    }

    // Clean up: delete the created user
    _, err = userGateway.DeleteUser(context.Background(), createdUser.UserID)
    if err != nil {
        t.Errorf("Failed to delete test user: %v", err)
    }
}

func TestGetUserHandler(t *testing.T) {
    ctx := context.Background()
    router := chi.NewRouter()
    testDBPool, err := pgxpool.New(ctx, "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    userGateway := db.NewUserTableGateway(testDBPool)
    sessionGateway := db.NewSessionTableGateway(testDBPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    router.Get("/api/v1/user/{userID}", handlers.HandleGetUser(authService))

    // Create a new user
    newUser := types.User{
        UserID: uuid.New(),
        Name:   "testuser",
    }

    _, err = userGateway.CreateUser(ctx, newUser)
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }

    // Create a new session for the user
    newSession := types.Session{
        ID:     uuid.New(),
        UserID: newUser.UserID,
    }

    _, err = sessionGateway.CreateSession(ctx, newSession)
    if err != nil {
        t.Fatalf("Failed to create test session: %v", err)
    }

    // Generate JWT token
    token, err := authService.GenerateJWT(ctx, newSession)
    if err != nil {
        t.Fatalf("Failed to generate JWT token: %v", err)
    }

    // Create request with JWT token in header
    req, err := http.NewRequest("GET", "/api/v1/user/"+newUser.UserID.String(), nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer "+token)

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var fetchedUser types.User
    err = json.NewDecoder(rr.Body).Decode(&fetchedUser)
    if err != nil {
        t.Errorf("Failed to decode response body: %v", err)
    }

    if fetchedUser.UserID != newUser.UserID {
        t.Errorf("Fetched user has incorrect UserID: got %v, want %v", fetchedUser.UserID, newUser.UserID)
    }

    if fetchedUser.Name != newUser.Name {
        t.Errorf("Fetched user has incorrect name: got %v, want %v", fetchedUser.Name, newUser.Name)
    }

    // Clean up: delete the created user and session
    _, err = userGateway.DeleteUser(ctx, newUser.UserID)
    if err != nil {
        t.Errorf("Failed to delete test user: %v", err)
    }

}

func TestValidateJWT(t *testing.T) {
    ctx := context.Background()
    testDBPool, err := pgxpool.New(ctx, "postgresql://myuser:mypassword@localhost:5432/goragdb")
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer testDBPool.Close()

    userGateway := db.NewUserTableGateway(testDBPool)
    sessionGateway := db.NewSessionTableGateway(testDBPool)
    authService := auth.NewAuthService(userGateway, sessionGateway)

    // Create a test user
    testUser := types.User{
        UserID: uuid.New(),
        Name:   "testuser",
    }



    // Insert the test user into the database
    success, err := userGateway.CreateUser(ctx, testUser)
    if err != nil || !success {
        t.Fatalf("Failed to create test user: %v", err)
    }
    defer func() {
        // Clean up: delete the test user
        _, err := userGateway.DeleteUser(ctx, testUser.UserID)
        if err != nil {
            t.Errorf("Failed to delete test user: %v", err)
        }
    }()

    // create a test session
    testSession := types.Session{
        ID:     uuid.New(),
        UserID: testUser.UserID,
    }

    // Insert the test session into the database
    success, err = sessionGateway.CreateSession(ctx, testSession)
    if err != nil || !success {
        t.Fatalf("Failed to create test session: %v", err)
    }

    // Generate JWT token for the test user
    token, err := authService.GenerateJWT(ctx, testSession)
    if err != nil {
        t.Fatalf("Failed to generate JWT: %v", err)
    }

    // Test ValidateJWT function
    validatedSession, err := authService.ValidateJWT(ctx, token)
    if err != nil {
        t.Errorf("ValidateJWT failed: %v", err)
    }

    if validatedSession.UserID != testUser.UserID {
        t.Errorf("Expected UserID %v, got %v", testUser.UserID, validatedSession.UserID)
    }

    if validatedSession.ID != testSession.ID {
        t.Errorf("Expected Name %v, got %v", testSession.ID, validatedSession.ID)
    }

    // Now delete the user from the database
    _, err = userGateway.DeleteUser(ctx, testUser.UserID)
    if err != nil {
        t.Fatalf("Failed to delete test user: %v", err)
    }

    // Try to validate the same token again
    _, err = authService.ValidateJWT(ctx, token)
    if err == nil {
        t.Errorf("Expected error when validating token for non-existent user, got nil")
    }
}
