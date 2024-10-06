package tests

import (
	"testing"
	// "github.com/stretchr/testify/assert"
	"rag-demo/pkg/index"
	"fmt"
)

func TestInitializeS3Client(t *testing.T) {

	// Test the initialization of the S3 client
	s3, err := index.NewS3Service()
	if err != nil {
		t.Fatalf("Failed to initialize S3 client: %v", err)
	}
	fmt.Printf("S3 client: %v\n", s3)
	// s3Client := InitializeS3Client()
	// assert.NotNil(t, s3Client, "S3 client should not be nil")
}