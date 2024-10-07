package index

import (
	"context"
	"fmt"
	"strings"
	"time"
	"rag-demo/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	textractTypes "github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/joho/godotenv"
	// "os"
	// "io"
)

type TextractService struct {
	Client *textract.Client
}


func NewTextractService() (*TextractService, error) {
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

	client := textract.NewFromConfig(cfg)

	return &TextractService{
		Client: client,
	}, nil
}

// StartTextDetection initiates the text detection process for a PDF
func (t *TextractService) StartTextDetection(ctx context.Context, bucket, documentKey string) (*textract.StartDocumentTextDetectionOutput, error) {
	input := &textract.StartDocumentTextDetectionInput{
		DocumentLocation: &textractTypes.DocumentLocation{
			S3Object: &textractTypes.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(documentKey),
			},
		},
		ClientRequestToken: aws.String("rag-demo"),
		JobTag:             aws.String("rag-demo"),
	}

	output, err := t.Client.StartDocumentTextDetection(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to start text detection: %w", err)
	}

	return output, nil
}


func (t *TextractService) GetTextFromPDF(ctx context.Context, jobID string, documentName string) (*types.DocumentText, error) {
	var fullText strings.Builder
	var nextToken *string

	for {
		input := &textract.GetDocumentTextDetectionInput{
			JobId:     aws.String(jobID),
			NextToken: nextToken,
		}

		output, err := t.Client.GetDocumentTextDetection(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to get text detection results: %w", err)
		}

		if output.JobStatus == textractTypes.JobStatusFailed {
			return nil, fmt.Errorf("text detection job failed: %s", aws.ToString(output.StatusMessage))
		}

		for _, block := range output.Blocks {
			if block.BlockType == textractTypes.BlockTypeLine {
				fullText.WriteString(aws.ToString(block.Text))
				fullText.WriteString(" ")
			}
		}

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken

		// Check if context has been cancelled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue processing
		}

		// Add a small delay to avoid hitting rate limits
		time.Sleep(100 * time.Millisecond)
	}

	// Split the full text into chunks of approximately one page each
	// Assuming an average of 3000 characters per page
	const charsPerPage = 3000
	text := fullText.String()
	var chunks []string
	for i := 0; i < len(text); i += charsPerPage {
		end := i + charsPerPage
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	return &types.DocumentText{
		Name:   documentName,
		Chunks: chunks,
	}, nil
}