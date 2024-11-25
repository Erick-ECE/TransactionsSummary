package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"transactions-summary/internal/infrastructure/database"
	"transactions-summary/internal/infrastructure/email"
	"transactions-summary/internal/infrastructure/file"
	"transactions-summary/internal/usecases"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	secretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	_ "github.com/go-sql-driver/mysql"
)

func getSecrets(ctx context.Context, secretName string, cfg aws.Config) (map[string]string, error) {
	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %v", err)
	}

	var secretMap map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &secretMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret string: %v", err)
	}

	return secretMap, nil
}

func handler(ctx context.Context, s3Event events.S3Event) {
	log.Println("Lambda function started processing S3 event")

	// Process each record in the S3 event
	for _, record := range s3Event.Records {
		s3Entity := record.S3
		bucketName := s3Entity.Bucket.Name
		objectKey := s3Entity.Object.Key

		log.Printf("Processing file: %s from bucket: %s", objectKey, bucketName)

		// Load AWS config
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Printf("Unable to load SDK config: %v", err)
			return
		}
		log.Println("AWS SDK config loaded successfully")

		// Retrieve secrets
		secretName := os.Getenv("SECRETS_MANAGER_NAME") // Set this in Lambda environment variables
		secrets, err := getSecrets(ctx, secretName, cfg)
		if err != nil {
			log.Printf("Failed to retrieve secrets: %v", err)
			return
		}
		log.Println("Secrets retrieved successfully from AWS Secrets Manager")

		// Access secrets
		dbUser := secrets["DB_USER"]
		dbPassword := secrets["DB_PASSWORD"]
		emailUser := secrets["EMAIL_USER"]
		emailPassword := secrets["EMAIL_PASSWORD"]

		// Load other environment variables
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPortStr := os.Getenv("SMTP_PORT")
		fromEmail := emailUser

		// Convert SMTP port from string to int
		smtpPort, err := strconv.Atoi(smtpPortStr)
		if err != nil {
			log.Printf("Invalid SMTP port: %v", err)
			continue
		}

		// Build the DSN (Data Source Name) for MySQL connection
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

		// Open the database connection
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Could not connect to the database: %v", err)
			continue
		}
		defer db.Close()

		// Test the database connection
		if err = db.Ping(); err != nil {
			log.Printf("Failed to ping database: %v", err)
			continue
		}
		log.Println("Successfully connected to the database")

		// Initialize repositories and services
		transactionRepo := database.NewMySQLTransactionRepo(db)
		csvReader := file.NewCSVReader()
		processTransactions := usecases.NewProcessTransactions(transactionRepo, csvReader)
		generateSummary := usecases.NewGenerateSummary(transactionRepo)
		emailService := email.NewGomailService(smtpHost, smtpPort, emailUser, emailPassword, fromEmail)
		sendSummaryEmail := usecases.NewSendSummaryEmail(generateSummary, emailService)

		// Create AWS S3 client using default credentials
		s3Client := s3.NewFromConfig(cfg)

		// Read the CSV file from S3
		log.Printf("Reading CSV file from S3: %s/%s", bucketName, objectKey)
		csvFileReader, err := readCSVFromS3(ctx, s3Client, bucketName, objectKey)
		if err != nil {
			log.Printf("Failed to read CSV from S3: %v", err)
			continue
		}
		log.Println("CSV file read successfully from S3")

		// Execute the ProcessTransactions use case
		accountToTransactions, err := processTransactions.Execute(csvFileReader)
		if err != nil {
			log.Printf("Could not process transactions: %v", err)
			continue
		}
		log.Println("Transactions processed successfully")

		// Send summary emails
		if err := sendSummaryEmail.Execute(accountToTransactions); err != nil {
			log.Printf("Could not send summary email: %v", err)
			continue
		}
		log.Println("Summary emails sent successfully")

		log.Printf("Successfully processed file: %s", objectKey)
	}
}

func main() {
	lambda.Start(handler)
}

// Helper function to read CSV from S3
func readCSVFromS3(ctx context.Context, client *s3.Client, bucket, key string) (*csv.Reader, error) {
	// Get the object from S3
	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	// Ensure the response body is closed after the reader finishes
	// Use io.ReadAll to fully read the content into memory first
	defer output.Body.Close()

	// Copy the response body to a buffer (to avoid closing issues)
	buffer, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %v", err)
	}

	// Create a new CSV reader from the buffer
	csvReader := csv.NewReader(bytes.NewReader(buffer))

	return csvReader, nil
}
