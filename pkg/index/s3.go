package index 

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/joho/godotenv"
	"rag-demo/types"
	"os"
	"io"
)

type S3Service struct {
	Client   *s3.Client
	Uploader *manager.Uploader
}


func NewS3Service() (*S3Service, error) {
	// Load environment variables from .env file if present
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}

	// Create a new S3 service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &S3Service{
		Client:   client,
		Uploader: uploader,
	}, nil
}

func (s *S3Service) CheckFileExists(ctx context.Context, bucket string, key string) (bool, error) {
	_, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s *S3Service) UploadFile(ctx context.Context, bucket string, key string, filename string) (types.Document, error) {

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return types.Document{}, err
	}

	output, err := s.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body: io.Reader(file),
	})

	if err != nil {
		return types.Document{}, err
	}
	res := types.Document{
		ObjectKey: *output.Key,
		FileName: filename,
	}


	return res, nil
}
