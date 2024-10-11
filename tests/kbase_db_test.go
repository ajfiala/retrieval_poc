package tests

import (
	"context"
	"rag-demo/pkg/db"
	"rag-demo/types"
	// "rag-demo/types"
	"fmt"

    "github.com/pgvector/pgvector-go"

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

func TestKbaseTableGateway(t *testing.T) {
    ctx := context.Background()

    pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }

    kbaseGateway := db.NewKbaseTableGateway(pool)

    // Test data
    testKbase := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test Kbase2",
        Description: "This is a test knowledge base",
    }

    // Register cleanup for testKbase before deferring pool.Close()
    defer func() {
        pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
        if err != nil {
            t.Fatalf("Failed to connect to database: %v", err)
        }
        _, err = pool.Exec(ctx, "DELETE FROM kbase WHERE uuid = $1", testKbase.ID)
        if err != nil {
            t.Logf("Failed to delete testKbase data: %v", err)
        }
        pool.Close()
    }()

    // Now defer pool.Close() so it's called last
    defer pool.Close()

    t.Run("CreateKbase", func(t *testing.T) {
        success, err := kbaseGateway.CreateKbase(ctx, testKbase)
        if err != nil {
            t.Fatalf("CreateKbase failed: %v", err)
        }
        if !success {
            t.Fatalf("CreateKbase returned false")
        }
    })

    t.Run("UpdateKbase", func(t *testing.T) {
        updatedTestKbase := types.Kbase{
            ID:          testKbase.ID,
            Name:        "Updated Test Kbase",
            Description: "This is an updated test knowledge base",
        }
        success, err := kbaseGateway.UpdateKbase(ctx, updatedTestKbase)
        if err != nil {
            t.Fatalf("UpdateKbase failed: %v", err)
        }
        if !success {
            t.Fatalf("UpdateKbase returned false")
        }
    })

    t.Run("GetKbase", func(t *testing.T) {
        kbase, err := kbaseGateway.GetKbase(ctx, testKbase.ID)
        if err != nil {
            t.Fatalf("GetKbase failed: %v", err)
        }
        if kbase.ID != testKbase.ID {
            t.Fatalf("GetKbase returned incorrect kbase")
        }
        if kbase.Name != "Updated Test Kbase" {
            t.Fatalf("GetKbase returned incorrect kbase")
        }
        if kbase.Description != "This is an updated test knowledge base" {
            t.Fatalf("GetKbase returned incorrect kbase")
        }
    })

    t.Run("ListKbases", func(t *testing.T) {
        // Test data. add another kbase
        testKbase2 := types.Kbase{
            ID:          uuid.New(),
            Name:        "Test Kbase123",
            Description: "This is a test knowledge base",
        }

        // Register cleanup for testKbase2 before pool.Close()
        defer func() {
            _, err = pool.Exec(ctx, "DELETE FROM kbase WHERE uuid = $1", testKbase2.ID)
            if err != nil {
                t.Logf("Failed to delete testKbase2 data: %v", err)
            }
        }()

        // Create another kbase
        success, err := kbaseGateway.CreateKbase(ctx, testKbase2)
        if err != nil {
            t.Fatalf("CreateKbase failed: %v", err)
        }
        if !success {
            t.Fatalf("CreateKbase returned false")
        }

        kbases, err := kbaseGateway.ListKbases(ctx)
        if err != nil {
            t.Fatalf("ListKbases failed: %v", err)
        }
        if len(kbases.Kbases) < 1 {
            t.Fatalf("ListKbases returned incorrect number of kbases")
        }
    })

	_, err = pool.Exec(ctx, "DELETE FROM kbase WHERE uuid = $1", testKbase.ID)
	if err != nil {
		t.Logf("Failed to delete testKbase data: %v", err)
	}
}

func TestKbaseEmbeddingTableGateway(t *testing.T) {
    ctx := context.Background()

    pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()

    kbaseGateway := db.NewKbaseTableGateway(pool)
    embeddingGateway := db.NewKbaseEmbeddingsTableGateway(pool)

    // Create a test Kbase
    testKbase := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test Kbase Embeddings",
        Description: "This is a test knowledge base for embeddings",
    }

    // Create the Kbase
    success, err := kbaseGateway.CreateKbase(ctx, testKbase)
    if err != nil || !success {
        t.Fatalf("Failed to create test Kbase: %v", err)
    }

    // Register cleanup for testKbase
    defer func() {
        _, err := pool.Exec(ctx, "DELETE FROM kbase WHERE uuid = $1", testKbase.ID)
        if err != nil {
            t.Logf("Failed to delete testKbase data: %v", err)
        }
    }()

    t.Run("CreateEmbedding", func(t *testing.T) {
        // Test data for embedding
        testEmbedding := types.KbaseEmbedding{
            UUID:      uuid.New(),
            KbaseID:   testKbase.ID,
            ChunkID:   1,
            Content:   "This is a test content chunk.",
			Embedding: pgvector.NewVector([]float32{1, 2, 3}),
            Metadata: map[string]interface{}{
                "source": "test",
            },
        }

        // Register cleanup for testEmbedding
        defer func() {
            _, err := pool.Exec(ctx, "DELETE FROM kbase_embeddings WHERE uuid = $1", testEmbedding.UUID)
            if err != nil {
                t.Logf("Failed to delete testEmbedding data: %v", err)
            }
        }()

        // Create the embedding
        success, err := embeddingGateway.CreateEmbedding(ctx, testEmbedding)
        if err != nil || !success {
            t.Fatalf("CreateEmbedding failed: %v", err)
        }

        // Optionally, retrieve and verify the embedding
        // (Implementation of GetEmbedding is needed)
    })
}
