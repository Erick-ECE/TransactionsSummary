package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"transactions-summary/internal/infrastructure/database"
	"transactions-summary/internal/infrastructure/email"
	"transactions-summary/internal/infrastructure/file"
	"transactions-summary/internal/usecases"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load database configuration (adjust as needed)
	dbUser := "erick"
	dbPassword := "3r1ck15"
	dbName := "transactions_summary"
	dbHost := "localhost"

	// Email configuration (adjust with your SMTP details)
	smtpHost := "smtp.gmail.com"
	smtpPort := 465
	emailUser := "#######"
	emailPassword := "##########"
	fromEmail := "#########"

	// Email recipient (could be passed as a parameter)
	recipientEmail := "ericken15@ciencias.unam.mx"

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
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
		log.Fatal("please provide the path to the CSV file")
	}
	filePath := os.Args[1]

	// Execute the ProcessTransactions use case
	err = processTransactions.Execute(filePath)
	if err != nil {
		log.Fatalf("could not process transactions: %v", err)
	}

	// Initialize the GenerateSummary use case
	generateSummary := usecases.NewGenerateSummary(transactionRepo)

	// Initialize the Gomail email service
	emailService := email.NewGomailService(smtpHost, smtpPort, emailUser, emailPassword, fromEmail)

	// Initialize and execute the SendSummaryEmail use case
	sendSummaryEmail := usecases.NewSendSummaryEmail(generateSummary, emailService)
	if err := sendSummaryEmail.Execute(recipientEmail); err != nil {
		log.Fatalf("could not send summary email: %v", err)
	}

	log.Println("Transactions processed and summary email sent successfully!")
}
