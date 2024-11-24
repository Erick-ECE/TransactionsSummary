package main

import (
	"fmt"
	"log"

	"transactions-summary/internal/infrastructure"
	"transactions-summary/internal/infrastructure/database"
	"transactions-summary/internal/infrastructure/file"
	"transactions-summary/internal/usecases"

	_ "github.com/go-sql-driver/mysql" // MySQL driver for database connection
	//"github.com/joho/godotenv"         // Package to load .env file
)

func main() {
	// Load environment variables
	infrastructure.LoadEnv()

	// Initialize the database
	db := infrastructure.InitDatabase()
	defer db.Close()

	// Initialize the database repository
	transactionRepo := database.NewMySQLTransactionRepo(db)

	// Initialize the CSV reader
	csvReader := file.NewCSVReader()

	// Initialize the email service
	emailService := infrastructure.InitEmailService()

	// Initialize AWS S3 configuration and client
	client, bucketName := infrastructure.InitAWSConfig()

	// Initialize the ProcessTransactions use case
	processTransactions := usecases.NewProcessTransactions(transactionRepo, csvReader)

	// Get the last uploaded file
	lastFile, err := infrastructure.GetLastUploadedFile(client, bucketName)
	if err != nil {
		log.Fatalf("Failed to find last uploaded file: %v", err)
	}

	fmt.Printf("Last uploaded file: %s\n", *lastFile.Key)

	// Read the file content
	FilecsvReader, err := infrastructure.ReadCSVFromS3(client, bucketName, *lastFile.Key)
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

	// Initialize and execute the SendSummaryEmail use case
	sendSummaryEmail := usecases.NewSendSummaryEmail(generateSummary, emailService)

	if err := sendSummaryEmail.Execute(accountToTransactions); err != nil {
		log.Fatalf("Could not send summary email: %v", err)
	}

	log.Println("Summary email sent successfully!")
}
