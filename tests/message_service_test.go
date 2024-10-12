package tests

import(
	"rag-demo/pkg/message"
	"rag-demo/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"rag-demo/types"
	"fmt"
	"testing"	
	"context"
	"github.com/google/uuid"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestMessageInitializeBedrock(t *testing.T) {
	// Test the initialization of the Textract client
	bed, err := message.NewBedrockRuntimeService("anthropic")
	if err != nil {
		fmt.Printf("Error initializing Bedrock Runtime client: %v\n", err)
		t.Fatalf("Error initializing Bedrock Runtime client: %v\n", err)
	}

	// assert that txt is of type TextractService
	assert.IsType(t, &message.BedrockRuntimeService{}, bed, "txt should be of type TextractService")
}

func TestInvokeModel(t *testing.T) {
	// create message
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

	// test messageRequest
	msg := types.MessageRequest{
		SessionId: testSession.ID,
		UserId: testUser.UserID,
		Text: "Hello, this is a test message",
	}

	// bed, err := message.NewBedrockRuntimeService("anthropic")
	// if err != nil {
	// 	fmt.Printf("Error initializing Bedrock Runtime client: %v\n", err)
	// 	t.Fatalf("Error initializing Bedrock Runtime client: %v\n", err)
	// }

	// // Invoke the model
	// err = bed.InvokeModel(ctx, msg)

	// // fmt.Printf("Invoke Bedrock output: %+v\n", output.GetStream().Reader)

	// assert.Nil(t, err, "Error should be nil")
	// assert.NotNil(t, output, "Output should not be nil")

	bed, err := message.NewBedrockRuntimeService("a121")
	if err != nil {
		fmt.Printf("Error initializing Bedrock Runtime client: %v\n", err)
		t.Fatalf("Error initializing Bedrock Runtime client: %v\n", err)
	}

	// Invoke the model
	err = bed.InvokeModel(ctx, msg)

	// fmt.Printf("Invoke Bedrock output: %+v\n", output.GetStream().Reader)

	assert.Nil(t, err, "Error should be nil")
}
