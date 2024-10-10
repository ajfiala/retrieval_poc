package tests

import (
	"testing"
	"rag-demo/pkg/index"
	"rag-demo/pkg/db"
	"fmt"
	"rag-demo/types"
	"rag-demo/pkg/embedding_orchestrator"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

)


func TestEmbeddingOrchestrator(t *testing.T) {
	// Test the UploadFile function
	txt, err := index.NewTextractService()
	if err != nil {
		t.Fatalf("Error initializing Textract client")
	}

	assert.NotNil(t, txt, "Textract client should not be nil")

	bucket := "lil-rag-kbase"
	key := "generative-ai-on-aws-how-to-choose.pdf"

	output, err := txt.StartTextDetection(context.Background(), bucket, key)

	fmt.Println("JobID: ", *output.JobId)

	// fmt.Printf("JobID: %s\n", *output.JobId)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, output, "Output should not be nil")

	// Get the results
	result, err := txt.GetTextFromPDF(context.Background(), *output.JobId, key)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, result, "Results should not be nil")

	// fmt.Println("Results: ", result)



	bedrockService, err := index.NewBedrockRuntimeService()
    if err != nil {
        t.Fatalf("Error initializing Bedrock service: %v", err)
    }

    // Assume you have a function to get your database pool
    dbPool := db.GetPool()
    dbService := db.NewKbaseEmbeddingsTableGateway(dbPool)
    kbaseDbService := db.NewKbaseTableGateway(dbPool)

	// Create a test kbase
	testKbase := types.Kbase{
		ID:          uuid.New(),
		Name:        "Test Kbase23",
		Description: "This is a test knowledge base",
	}

	// Register cleanup for testKbase and its embeddings
	defer func() {
		// First, delete from kbase_embeddings
		_, err := dbPool.Exec(context.Background(), "DELETE FROM kbase_embeddings WHERE kbase_id = $1", testKbase.ID)
		if err != nil {
			t.Logf("Failed to delete kbase_embeddings data: %v", err)
		}

		// Then, delete from kbase
		_, err = dbPool.Exec(context.Background(), "DELETE FROM kbase WHERE uuid = $1", testKbase.ID)
		if err != nil {
			t.Logf("Failed to delete kbase data: %v", err)
		}
	}()

	kbaseDbService.CreateKbase(context.Background(), testKbase)


    // Create the orchestrator
    orchestrator := orchestrator.NewOrchestrator(bedrockService, dbService)

	err = orchestrator.ProcessAndStoreEmbeddings(context.Background(), *result, testKbase.ID)
	assert.Nil(t, err, "Error should be nil")
	if err != nil {
		fmt.Println(err)
	}


	// fmt.Printf("Results: %+v\n", result)
}