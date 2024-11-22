package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver for database connection

	"transactions-summary/internal/infrastructure/database"
	"transactions-summary/internal/infrastructure/file"
	"transactions-summary/internal/usecases"
)

func main() {
	// Load database configuration (you might want to use environment variables or a config file)
	dbUser := "####"
	dbPassword := "####"
	dbName := "transactions_summary"
	dbHost := "localhost" // Change this if needed

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

	log.Println("Transactions processed successfully!")
}
