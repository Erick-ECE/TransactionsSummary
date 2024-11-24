package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"transactions-summary/internal/infrastructure/database"
	"transactions-summary/internal/infrastructure/email"
	"transactions-summary/internal/infrastructure/file"
	"transactions-summary/internal/usecases"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	_ "github.com/go-sql-driver/mysql" // MySQL driver for database connection
	"github.com/joho/godotenv"         // Package to load .env file
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	// Load email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT") // Port is loaded as a string
	emailUser := os.Getenv("EMAIL_USER")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	fromEmail := os.Getenv("EMAIL_USER") // Usually the same as emailUser

	// Convert SMTP port from string to int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Invalid SMTP port: %v", err)
	}

	// Build the DSN (Data Source Name) for MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	// Initialize the database repository
	transactionRepo := database.NewMySQLTransactionRepo(db)

	// Initialize the CSV reader
	csvReader := file.NewCSVReader()

	// Initialize the ProcessTransactions use case
	processTransactions := usecases.NewProcessTransactions(transactionRepo, csvReader)

	// Get CSV file from S#
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

	// Get the last uploaded file
	lastFile, err := getLastUploadedFile(client, bucketName)
	if err != nil {
		log.Fatalf("Failed to find last uploaded file: %v", err)
	}

	fmt.Printf("Last uploaded file: %s\n", *lastFile.Key)

	// Read the file content
	FilecsvReader, err := readCSVFromS3(client, bucketName, *lastFile.Key)
	if err != nil {
		log.Fatalf("Failed to read CSV from S3: %v", err)
	}

	// Execute the ProcessTransactions use case
	accountToTransactions, err := processTransactions.Execute(FilecsvReader)
	if err != nil {
		log.Fatalf("Could not process transactions: %v", err)
	}

	// Initialize the GenerateSummary use case
	generateSummary := usecases.NewGenerateSummary(transactionRepo)

	// Initialize the Gomail email service
	emailService := email.NewGomailService(smtpHost, smtpPort, emailUser, emailPassword, fromEmail)

	// ============================================

	// Initialize and execute the SendSummaryEmail use case
	sendSummaryEmail := usecases.NewSendSummaryEmail(generateSummary, emailService)

	if err := sendSummaryEmail.Execute(accountToTransactions); err != nil {
		log.Fatalf("Could not send summary email: %v", err)
	}

	log.Println("Summary email sent successfully!")
}

// getLastUploadedFile fetches the most recently uploaded file in the given bucket path.
func getLastUploadedFile(client *s3.Client, bucket string) (*types.Object, error) {
	// List objects in the bucket with the specified prefix
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

// readFileFromS3 reads a file from S3 and returns its content as a string.
func readCSVFromS3(client *s3.Client, bucket, key string) (*csv.Reader, error) {
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
