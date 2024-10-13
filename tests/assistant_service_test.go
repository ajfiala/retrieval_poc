package tests

import (
	"context"
	"rag-demo/pkg/db"
	"rag-demo/types"
	"rag-demo/pkg/assistant"
	// "rag-demo/types"
	"fmt"
	"testing"
	"sync"
	"github.com/google/uuid"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

func TestInitializeAssistantService(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

	assistantGateway := db.NewAssistantTableGateway(pool)
	assistantService := assistant.NewAssistantService(assistantGateway)
	// assert that assistantService is of type AssistantService
	assert.IsType(t, &assistant.AssistantServiceImpl{}, assistantService, "assistantService should be of type AssistantService")
}

func TestCreateAssistant(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	assistantGateway := db.NewAssistantTableGateway(pool)
	assistantService := assistant.NewAssistantService(assistantGateway)

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

	// make waitgroup and result channel
	resultCh := make(types.ResultChannel, 1) 
	wg := &sync.WaitGroup{}
	wg.Add(1)


	assistantService.CreateAssistant(ctx, testAssistant, resultCh, wg)

	wg.Wait()

	result := <-resultCh

	fmt.Println("Result: ", result)

	assert.True(t, result.Success, "Result should be successful")
	
	assert.NoError(t, result.Error, "CreateAssistant should not return an error")

	// clean up by deleting the assistant. use exec 
	_, err = pool.Exec(ctx, "DELETE FROM assistant WHERE uuid = $1", testAssistant.ID)
	if err != nil {
		t.Fatalf("Error deleting assistant: %v", err)
	}

}

func TestListAssistants(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	assistantGateway := db.NewAssistantTableGateway(pool)
	assistantService := assistant.NewAssistantService(assistantGateway)

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
		Name:     "insurance assistant23",
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

	success, err := assistantGateway.CreateAssistant(ctx, testAssistant)
	if err != nil {
		t.Fatalf("CreateAssistant failed: %v", err)
	}
	if !success {
		t.Fatalf("CreateAssistant returned false")
	}
	success, err = assistantGateway.CreateAssistant(ctx, testAssistant2)
	if err != nil {
		t.Fatalf("CreateAssistant failed: %v", err)
	}
	if !success {
		t.Fatalf("CreateAssistant returned false")
	}

	// make waitgroup and result channel
	resultCh := make(types.ResultChannel, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)

	assistantService.ListAssistants(ctx, resultCh, wg)

	wg.Wait()

	result := <-resultCh

	fmt.Println("Result: ", result)
	assistantList := result.Data

	assert.True(t, result.Success, "Result should be successful")
	assert.IsType(t, types.AssistantList{}, assistantList, "Result data should be of type AssistantList")

	// clean up by deleting the assistants. use exec
	_, err = pool.Exec(ctx, "DELETE FROM assistant")
	if err != nil {
		t.Fatalf("Error deleting assistant: %v", err)
	}
}

