package index 

import (
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"context"
	"github.com/joho/godotenv"
	"encoding/json"
	"rag-demo/types"
	"fmt"
)

type BedrockRuntimeService struct {
	Client *bedrockruntime.Client
}

func NewBedrockRuntimeService() (*BedrockRuntimeService, error) {
	// Load environment variables from .env file if present
	if err := godotenv.Load("../.env"); err != nil {
		// Handle error if .env file is not found
		// For testing, we can set default environment variables
		// fmt.Printf("Error loading .env file")
	}

	// Create a new S3 service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := bedrockruntime.NewFromConfig(cfg)

	return &BedrockRuntimeService{
		Client: client,
	}, nil
}

func (b *BedrockRuntimeService) GetEmbeddings(ctx context.Context, doctext types.DocumentText) (*bedrockruntime.InvokeModelOutput, error) {
	// Prepare the input for the Titan embedding model
	inputStruct := types.TitanEmbeddingInput{
		InputText: doctext.Chunks[0], // Using the first chunk as an example
	}

	// Convert the input struct to JSON
	inputJSON, err := json.Marshal(inputStruct)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input: %w", err)
	}

	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String("amazon.titan-embed-g1-text-02"),
		Body:        inputJSON,
		ContentType: aws.String("application/json"),
	}

	output, err := b.Client.InvokeModel(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error invoking model: %w", err)
	}


	return output, nil
}


