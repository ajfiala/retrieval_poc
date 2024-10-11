package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	
	"context"

	"rag-demo/pkg/index"

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
