package tests

import (
	"context"
	"rag-demo/pkg/db"
	"rag-demo/types"
	// "rag-demo/types"
	"fmt"
	"testing"
	"github.com/google/uuid"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
    // Load environment variables from .env file if present
    if err := godotenv.Load("../.env"); err != nil {
        // Handle error if .env file is not found
        // For testing, we can set default environment variables
		fmt.Printf("Error loading .env file")
    }
}

func TestAssistantTableGateway(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

	assistantGateway := db.NewAssistantTableGateway(pool)

    // Test data


	// Test data
	testAssistant := types.Assistant{
		ID:       uuid.New(),
		Name:     "insurance assistant3",
		Type:     "rag",
		Model:    "anthropic.claude-3-5-sonnet-20240620-v1:0",
		SystemPrompts: 
`You are an insurance manager for a large Thai insurance company.
You are an expert on insurance products, rules, services, and policies. Users will
ask you questions about insurance products, rules, services, and policies. Respond
using the provided context. If the question is general enough, you can provide a
general answer. Otherwise, don't answer if you aren't sure or if the answer cannot be found
in the provided context. In this case, direct the user to another source of information.`,
		Metadata: &types.Metadata{
			Title:       "Insurance Assistant",
			Description: "Assists users with insurance-related queries.",
			Icon:        "insurance_icon.png",
			Prompts:     []string{"How can I assist you with insurance today?",
		"How do I file a claim?",
		"What is the process for renewing my policy?",},
		},
	}

	t.Run("CreateAssistant", func(t *testing.T) {
        success, err := assistantGateway.CreateAssistant(ctx, testAssistant)
        if err != nil {
            t.Fatalf("CreateAssistant failed: %v", err)
        }
        if !success {
            t.Fatalf("CreateAssistant returned false")
        }
    })

	t.Run("GetAssistant", func(t *testing.T) {
		assistant, err := assistantGateway.GetAssistant(ctx, testAssistant.ID)
		if err != nil {
			t.Fatalf("GetAssistant failed: %v", err)
		}
		if assistant.ID != testAssistant.ID {
			t.Errorf("Expected ID %v, got %v", testAssistant.ID, assistant.ID)
		}
		if assistant.Name != testAssistant.Name {
			t.Errorf("Expected Name %v, got %v", testAssistant.Name, assistant.Name)
		}
		if assistant.Type != testAssistant.Type {
			t.Errorf("Expected Type %v, got %v", testAssistant.Type, assistant.Type)
		}
		if assistant.Model != testAssistant.Model {
			t.Errorf("Expected Model %v, got %v", testAssistant.Model, assistant.Model)
		}
		if assistant.SystemPrompts != testAssistant.SystemPrompts {
			t.Errorf("Expected SystemPrompts %v, got %v", testAssistant.SystemPrompts, assistant.SystemPrompts)
		}
		if assistant.Metadata.Title != testAssistant.Metadata.Title {
			t.Errorf("Expected Metadata.Title %v, got %v", testAssistant.Metadata.Title, assistant.Metadata.Title)
		}
		if assistant.Metadata.Description != testAssistant.Metadata.Description {
			t.Errorf("Expected Metadata.Description %v, got %v", testAssistant.Metadata.Description, assistant.Metadata.Description)
		}
		if assistant.Metadata.Icon != testAssistant.Metadata.Icon {
			t.Errorf("Expected Metadata.Icon %v, got %v", testAssistant.Metadata.Icon, assistant.Metadata.Icon)
		}
	})

	t.Run("ListAssistants", func(t *testing.T) {
		assistants, err := assistantGateway.ListAssistants(ctx)
		if err != nil {
			t.Fatalf("ListAssistants failed: %v", err)
		}
		fmt.Println(assistants)
		if len(assistants.Assistants) == 0 {
			t.Fatalf("ListAssistants returned empty list")
		}
	})

	t.Run("UpdateAssistant", func(t *testing.T) {
		updatedAssistant := testAssistant
		updatedAssistant.Name = "updated insurance assistant"

		success, err := assistantGateway.UpdateAssistant(ctx, updatedAssistant)
		if err != nil {
			t.Fatalf("UpdateAssistant failed: %v", err)
		}
		if !success {
			t.Fatalf("UpdateAssistant returned false")
		}
	})

	t.Run("DeleteAssistant", func(t *testing.T) {
		_, err := pool.Exec(ctx, "DELETE FROM assistant WHERE uuid = $1", testAssistant.ID)
		if err != nil {
			t.Fatalf("DeleteAssistant failed: %v", err)
		}
	})

}