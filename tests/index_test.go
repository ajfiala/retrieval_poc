package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"context"
	// "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"rag-demo/pkg/index"
	"rag-demo/pkg/db"
	"rag-demo/pkg/kbase"
	"github.com/lib/pq"
	"github.com/google/uuid"
	"sync"
	"rag-demo/types"
	"fmt"
)

func TestInitializeS3Client(t *testing.T) {

	// Test the initialization of the S3 client
	s3, err := index.NewS3Service()
	if err != nil {
		fmt.Printf("Error initializing S3 client: %v\n", err)
		t.Fatalf("Error initializing	S3 client: %v\n", err)
	}

	// assert that s3 is of typ S3Service
	assert.IsType(t, &index.S3Service{}, s3, "s3 should be of type S3Service")
}

func TestCheckFileExists(t *testing.T) {
	// Test the CheckFileExists function
	s3, err := index.NewS3Service()
	if err != nil {
		t.Fatalf("Error initializing S3 client")
	}
		
	assert.NotNil(t, s3, "S3 client should not be nil")

	bucket := "lil-rag-kbase"
	key := "visa_types_and_fees.pdf"
	exists, err := s3.CheckFileExists(context.Background(), bucket, key)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, exists, "File should exist")
}

func TestUploadFile(t *testing.T) {
	// Test the UploadFile function
	s3, err := index.NewS3Service()
	if err != nil {
		t.Fatalf("Error initializing S3 client")
	}

	assert.NotNil(t, s3, "S3 client should not be nil")

	bucket := "lil-rag-kbase"
	key := "browns_letter_1974.pdf"
	filename := "../browns_letter_1974.pdf"
	output, err := s3.UploadFile(context.Background(), bucket, key, filename)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, output, "Output should not be nil")

	// assert output

}

func TestInitializeTextract(t *testing.T) {
	// Test the initialization of the Textract client
	txt, err := index.NewTextractService()
	if err != nil {
		fmt.Printf("Error initializing Textract client: %v\n", err)
		t.Fatalf("Error initializing Textract client: %v\n", err)
	}

	// assert that txt is of type TextractService
	assert.IsType(t, &index.TextractService{}, txt, "txt should be of type TextractService")
}

func TestStartTextract(t *testing.T) {
	// Test the UploadFile function
	txt, err := index.NewTextractService()
	if err != nil {
		t.Fatalf("Error initializing Textract client")
	}

	assert.NotNil(t, txt, "Textract client should not be nil")

	bucket := "lil-rag-kbase"
	key := "browns_letter_1974.pdf"

	output, err := txt.StartTextDetection(context.Background(), bucket, key)

	// fmt.Printf("JobID: %s\n", *output.JobId)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, output, "Output should not be nil")
}

func TestGetTextractResults(t *testing.T) {
	// Test the UploadFile function
	txt, err := index.NewTextractService()
	if err != nil {
		t.Fatalf("Error initializing Textract client")
	}

	assert.NotNil(t, txt, "Textract client should not be nil")

	bucket := "lil-rag-kbase"
	key := "browns_letter_1974.pdf"

	output, err := txt.StartTextDetection(context.Background(), bucket, key)

	// fmt.Printf("JobID: %s\n", *output.JobId)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, output, "Output should not be nil")

	// Get the results
	result, err := txt.GetTextFromPDF(context.Background(), *output.JobId, key)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, result, "Results should not be nil")

	// fmt.Printf("Results: %+v\n", result)
}

func TestInitializeBedrock(t *testing.T) {
	// Test the initialization of the Textract client
	bed, err := index.NewBedrockRuntimeService()
	if err != nil {
		fmt.Printf("Error initializing Bedrock Runtime client: %v\n", err)
		t.Fatalf("Error initializing Bedrock Runtime client: %v\n", err)
	}

	// assert that txt is of type TextractService
	assert.IsType(t, &index.BedrockRuntimeService{}, bed, "txt should be of type TextractService")
}

func TestGetEmbeddings(t *testing.T) {
	// Test the UploadFile function
	txt, err := index.NewTextractService()
	if err != nil {
		t.Fatalf("Error initializing Textract client")
	}

	assert.NotNil(t, txt, "Textract client should not be nil")

	bucket := "lil-rag-kbase"
	key := "browns_letter_1974.pdf"

	output, err := txt.StartTextDetection(context.Background(), bucket, key)

	// fmt.Printf("JobID: %s\n", *output.JobId)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, output, "Output should not be nil")

	// Get the results
	result, err := txt.GetTextFromPDF(context.Background(), *output.JobId, key)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, result, "Results should not be nil")

	// fmt.Printf("Results: %+v\n", result)

	bed, err := index.NewBedrockRuntimeService()
	if err != nil {
		fmt.Printf("Error initializing Bedrock Runtime client: %v\n", err)
		t.Fatalf("Error initializing Bedrock Runtime client: %v\n", err)
	}

	// Get the embeddings
	embeddings, err := bed.GetEmbeddings(context.Background(), *result)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, embeddings, "Embeddings should not be nil")

	// assert that embeddings.Body is []float32
	assert.IsType(t, []uint8{}, embeddings.Body, "embeddings.Body should be of type []float32")

	// fmt.Printf("Invoke Bedrock Embeddings output: %+v\n", embeddings)
}

func TestStoreEmbeddings(t *testing.T) {
	// Test the UploadFile function
	txt, err := index.NewTextractService()
	if err != nil {
		t.Fatalf("Error initializing Textract client")
	}

	assert.NotNil(t, txt, "Textract client should not be nil")

	bucket := "lil-rag-kbase"
	key := "browns_letter_1974.pdf"

	output, err := txt.StartTextDetection(context.Background(), bucket, key)

	fmt.Printf("JobID: %s\n", *output.JobId)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, output, "Output should not be nil")

	// Get the results
	result, err := txt.GetTextFromPDF(context.Background(), *output.JobId, key)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, result, "Results should not be nil")

	// fmt.Printf("Results: %+v\n", result)

	bed, err := index.NewBedrockRuntimeService()
	if err != nil {
		fmt.Printf("Error initializing Bedrock Runtime client: %v\n", err)
		t.Fatalf("Error initializing Bedrock Runtime client: %v\n", err)
	}

	// Get the embeddings
	embeddings, err := bed.GetEmbeddings(context.Background(), *result)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, embeddings, "Embeddings should not be nil")

	// assert that embeddings.Body is []float32
	assert.IsType(t, []uint8{}, embeddings.Body, "embeddings.Body should be of type []float32")

	// fmt.Printf("Invoke Bedrock Embeddings output: %+v\n", embeddings)

	// create kbase 
	ctx := context.Background()
    dbPool := setupTestKbaseDB(t)
    defer dbPool.Close()

    kbaseEmbeddingsGateway := db.NewKbaseEmbeddingsTableGateway(dbPool)
	kbaseGateway := db.NewKbaseTableGateway(dbPool)
    kbaseService := kbase.NewKbaseService(kbaseGateway)

    // Test data
    testKbase := types.Kbase{
        ID:          uuid.New(),
        Name:        "Test Kbase999",
        Description: "Test description for Kbase",
    }

    // Create a new kbase
    resultCh := make(types.ResultChannel, 1) // Buffered channel to prevent deadlock
    wg := &sync.WaitGroup{}
    wg.Add(1)

    go kbaseService.CreateKbase(ctx, testKbase, resultCh, wg)

    wg.Wait()             // Wait for the goroutine to finish
    kbaseResult := <-resultCh  // Read the result from the channel

    fmt.Println("kbaseResult: ", kbaseResult)
    assert.True(t, kbaseResult.Success, "Result should be successful")
    assert.NoError(t, kbaseResult.Error, "CreateKbase should not return an error")

	// Define a struct to match the Bedrock response structure
	type BedrockEmbeddingResponse struct {
		Embedding []float64 `json:"embedding"`
	}

	// Unmarshal the Body into our struct
	var embeddingResponse BedrockEmbeddingResponse
	err = json.Unmarshal(embeddings.Body, &embeddingResponse)
	assert.Nil(t, err, "Error unmarshaling embeddings should be nil")

	// Print out the structure
	fmt.Printf("Embeddings: %+v\n", embeddingResponse.Embedding)

	// Convert []float64 to pq.Float64Array
	pqEmbedArray := pq.Float64Array(embeddingResponse.Embedding)

	// Test data for embedding
	testEmbedding := types.KbaseEmbedding{
		UUID:      uuid.New(),
		KbaseID:   testKbase.ID,
		ChunkID:   1,
		Content:   result.Chunks[0],
		Embedding: pqEmbedArray,
		Metadata: map[string]interface{}{
			"source": "test",
		},
	}

	// Store the embeddings
	embeddingRes, err := kbaseEmbeddingsGateway.CreateEmbedding(context.Background(), testEmbedding)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, embeddingRes, "Embedding should be stored successfully")
}