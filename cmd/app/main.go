package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"transactions-summary/internal/infrastructure/database"
	"transactions-summary/internal/infrastructure/email"
	"transactions-summary/internal/infrastructure/file"
	"transactions-summary/internal/usecases"

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

	// Email recipient (hardcoded or from an environment variable)
	recipientEmail := "ericken15@ciencias.unam.mx"

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

	// Get the CSV file path from the command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the CSV file")
	}
	filePath := os.Args[1]

	// Execute the ProcessTransactions use case
	log.Printf("Processing transactions from CSV file: %s", filePath)
	err = processTransactions.Execute(filePath)
	if err != nil {
		log.Fatalf("Could not process transactions: %v", err)
	}
	log.Println("Transactions successfully processed")

	// Initialize the GenerateSummary use case
	generateSummary := usecases.NewGenerateSummary(transactionRepo)

	// Initialize the Gomail email service
	emailService := email.NewGomailService(smtpHost, smtpPort, emailUser, emailPassword, fromEmail)

	// Initialize and execute the SendSummaryEmail use case
	sendSummaryEmail := usecases.NewSendSummaryEmail(generateSummary, emailService)
	log.Printf("Sending summary email to: %s", recipientEmail)
	if err := sendSummaryEmail.Execute(recipientEmail); err != nil {
		log.Fatalf("Could not send summary email: %v", err)
	}

	log.Println("Summary email sent successfully!")
}
