package message

import (
	"rag-demo/types"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/aws"
	"context"
	"os"
	"github.com/joho/godotenv"
	"fmt"
)

type BedrockRuntimeService struct {
	Client *bedrockruntime.Client
	Provider Provider
}

func NewBedrockRuntimeService(providerName string) (*BedrockRuntimeService, error) {
	if err := godotenv.Load("../.env"); err != nil {
		// Handle error or continue if .env file doesn't exist
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := bedrockruntime.NewFromConfig(cfg)

	provider, err := GetProvider(providerName)
	if err != nil {
		return nil, err
	}

	return &BedrockRuntimeService{
		Client:   client,
		Provider: provider,
	}, nil
}

func (b *BedrockRuntimeService) InvokeModel(ctx context.Context, message types.MessageRequest) (string, error) {
	jsonData, err := b.Provider.BuildRequest(message)
	if err != nil {
		return "", fmt.Errorf("error building request: %w", err)
	}

	modelID := os.Getenv("BEDROCK_MODEL_ID") // Ensure this is set per provider
	input := &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(modelID),
		Body:        jsonData,
		ContentType: aws.String("application/json"),
	}

	output, err := b.Client.InvokeModelWithResponseStream(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error invoking model: %w", err)
	}

	response, err := b.Provider.ProcessResponse(output)
	if err != nil {
		return "", fmt.Errorf("error processing response: %w", err)
	}

	return response, nil
}