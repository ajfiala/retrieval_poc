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
	"github.com/google/uuid"
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
    // Generate a unique client request token
    clientRequestToken := uuid.New().String()

    input := &textract.StartDocumentTextDetectionInput{
        DocumentLocation: &textractTypes.DocumentLocation{
            S3Object: &textractTypes.S3Object{
                Bucket: aws.String(bucket),
                Name:   aws.String(documentKey),
            },
        },
        ClientRequestToken: aws.String(clientRequestToken),
        JobTag:             aws.String("rag-demo-" + documentKey), // Make job tag unique per document
    }

    output, err := t.Client.StartDocumentTextDetection(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to start text detection: %w", err)
    }

    return output, nil
}


func (t *TextractService) GetTextFromPDF(ctx context.Context, jobID string, documentName string) (*types.DocumentText, error) {
    var fullText strings.Builder
    var blockCount int

    // Poll for job completion
    for {
        status, err := t.getJobStatus(ctx, jobID)
        if err != nil {
            return nil, fmt.Errorf("failed to get job status: %w", err)
        }

        fmt.Printf("Job Status: %s\n", status)

        if status == textractTypes.JobStatusSucceeded {
            break
        }

        if status == textractTypes.JobStatusFailed {
            return nil, fmt.Errorf("text detection job failed")
        }

        // Wait before checking again
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(5 * time.Second):
            // Continue polling
        }
    }

    // Process results
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

        fmt.Printf("Number of blocks in this response: %d\n", len(output.Blocks))

        for _, block := range output.Blocks {
            if block.BlockType == textractTypes.BlockTypeLine {
                fullText.WriteString(aws.ToString(block.Text))
                fullText.WriteString(" ")
                blockCount++
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
    }

    fmt.Printf("Total number of text blocks processed: %d\n", blockCount)
    fmt.Printf("Total text length: %d characters\n", fullText.Len())

    // Split the full text into chunks of approximately one page each
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

    fmt.Printf("Number of chunks created: %d\n", len(chunks))

    return &types.DocumentText{
        Name:   documentName,
        Chunks: chunks,
    }, nil
}

// Helper function to get job status
func (t *TextractService) getJobStatus(ctx context.Context, jobID string) (textractTypes.JobStatus, error) {
    input := &textract.GetDocumentTextDetectionInput{
        JobId: aws.String(jobID),
    }

    output, err := t.Client.GetDocumentTextDetection(ctx, input)
    if err != nil {
        return "", fmt.Errorf("failed to get job status: %w", err)
    }

    return output.JobStatus, nil
}