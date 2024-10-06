package index

import (
	"context"
	// "fmt"
	// "os"
	"github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/aws/credentials"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/joho/godotenv"
)

type S3Service struct {
	Uploader *manager.Uploader
}

type S3ServiceInterface interface {
	UploadFile(ctx context.Context, bucket string, key string, file string) error
	CheckFileExists(ctx context.Context, bucket string, key string) (bool, error)
}

func NewS3Service() S3ServiceInterface {
	// Load environment variables from .env file if present
	if err := godotenv.Load("../.env"); err != nil {
		// Handle error if .env file is not found
		// For testing, we can set default environment variables
		// fmt.Printf("Error loading .env file")
	}

	// Create a new S3 service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("failed to load configuration, " + err.Error())
	}

	uploader := s3manager.NewUploaderFromConfig(cfg)

	return &S3Service{Uploader: uploader}
}


// Upload uploads a file to the configured S3 bucket.
// It takes the local filename, reads the file, and uploads it to S3.
// Returns the S3 object key or an error.
func (s *s3ServiceImpl) Upload(ctx context.Context, filename string) (string, error) {
    // Open the file for reading
    file, err := os.Open(filename)
    if err != nil {
        return "", fmt.Errorf("failed to open file %q: %v", filename, err)
    }
    defer file.Close()

    // Generate a unique object key
    uniqueID := filepath.Base(filename) // You can modify this to add UUIDs if needed
    key := uniqueID

    // Upload the file to S3
    result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(key),
        Body:   file,
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload file to S3: %v", err)
    }

    return result.Location, nil
}