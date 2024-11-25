package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Check if filename is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_presigned_url.go <filename>")
		os.Exit(1)
	}
	key := os.Args[1]

	bucket := "transactions-demo" // Update with your bucket name
	region := "us-east-2"         // Update with your AWS region

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create a presigner
	presigner := s3.NewPresignClient(client)

	params := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	presignedURL, err := presigner.PresignPutObject(context.TODO(), params, func(po *s3.PresignOptions) {
		po.Expires = 15 * time.Minute
	})
	if err != nil {
		log.Fatalf("Failed to generate presigned URL: %v", err)
	}

	fmt.Println("Upload the file using this URL:")
	fmt.Println(presignedURL.URL)
	fmt.Printf("File will be uploaded as: %s\n", key)
}
