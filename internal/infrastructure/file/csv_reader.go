package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"transactions-summary/internal/entities"
)

// CSVReader implements the FileReader interface to read transactions from a CSV file.
type CSVReader struct{}

// NewCSVReader creates a new CSVReader instance.
func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadTransactions reads a CSV file and returns a list of transactions.
func (r *CSVReader) ReadTransactions(filePath string) ([]entities.Transaction, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV: %v", err)
	}

	var transactions []entities.Transaction

	// Skip the first row (header)
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		// Example CSV structure: Id,Date,Transaction
		// record[0] -> Id, record[1] -> Date (MM/DD), record[2] -> Transaction Amount (signed)

		// Parse the transaction amount
		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction amount in CSV: %v", err)
		}

		// Parse the date (assuming the current year)
		monthDay := strings.Split(record[1], "/")
		if len(monthDay) != 2 {
			return nil, fmt.Errorf("invalid date format in CSV: %s", record[1])
		}
		month, _ := strconv.Atoi(monthDay[0])
		day, _ := strconv.Atoi(monthDay[1])

		// Construct the date with the current year
		year := time.Now().Year()
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		// Create a transaction object
		txn := entities.Transaction{
			ID:     record[0],
			Amount: amount,
			Date:   date,
			Type:   determineTransactionType(amount), // Determine type based on the sign of the amount
		}

		transactions = append(transactions, txn)
	}

	return transactions, nil
}

// determineTransactionType returns "debit" if the amount is negative, otherwise "credit".
func determineTransactionType(amount float64) string {
	if amount < 0 {
		return "debit"
	}
	return "credit"
}
