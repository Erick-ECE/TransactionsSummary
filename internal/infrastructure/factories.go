package infrastructure

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"transactions-summary/internal/infrastructure/email"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from the .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// InitDatabase initializes the database connection
func InitDatabase() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	return db
}

// InitEmailService initializes the email service
func InitEmailService() *email.GomailService {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	emailUser := os.Getenv("EMAIL_USER")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	// Convert SMTP port from string to int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Invalid SMTP port: %v", err)
	}

	return email.NewGomailService(smtpHost, smtpPort, emailUser, emailPassword, emailUser)
}

// InitAWSConfig initializes the AWS S3 configuration and client
func InitAWSConfig() (*s3.Client, string) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("S3_BUCKET_NAME")

	// Create a custom AWS config
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      region,
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	return client, bucketName
}

// GetLastUploadedFile fetches the most recently uploaded file in the given bucket
func GetLastUploadedFile(client *s3.Client, bucket string) (*types.Object, error) {
	// List objects in the bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	if len(output.Contents) == 0 {
		return nil, fmt.Errorf("no files found in bucket path: %s", bucket)
	}

	// Find the most recently modified object
	var lastObject *types.Object
	for _, obj := range output.Contents {
		if lastObject == nil || obj.LastModified.After(*lastObject.LastModified) {
			temp := obj
			lastObject = &temp
		}
	}

	return lastObject, nil
}

// ReadCSVFromS3 reads a file from S3 and returns its content as a CSV reader
func ReadCSVFromS3(client *s3.Client, bucket, key string) (*csv.Reader, error) {
	// Get the object from S3
	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	// Return a CSV reader
	return csv.NewReader(output.Body), nil
}
