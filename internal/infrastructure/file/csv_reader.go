package file

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"transactions-summary/internal/entities"
)

// CSVReader implements the FileReader interface to read transactions from a CSV file.
type CSVReader struct{}

// NewCSVReader creates a new CSVReader instance.
func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadTransactions reads a CSV file and returns a list of transactions.
func (r *CSVReader) ReadTransactions(reader *csv.Reader) ([]entities.Transaction, error) {

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

		// Example CSV structure: Id,Date,Transaction,AccountId

		// Parse AccountId (last column)
		accountId := strings.TrimSpace(record[2])

		// Parse the transaction amount
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction amount in CSV: %v", err)
		}

		// Parse the date (assuming the current year)
		monthDay := strings.Split(record[0], "/")
		if len(monthDay) != 2 {
			return nil, fmt.Errorf("invalid date format in CSV: %s", record[0])
		}
		month, _ := strconv.Atoi(monthDay[0])
		day, _ := strconv.Atoi(monthDay[1])

		// Construct the date with the current year
		year := time.Now().Year()
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		// Determine transaction type using the function
		transactionType := determineTransactionType(amount)

		newUUID := uuid.New()
		// Create a transaction object
		transaction := entities.Transaction{
			ID:              newUUID.String(),
			AccountID:       accountId,
			Amount:          amount,
			TransactionDate: date,
			Type:            transactionType,
		}
		transactions = append(transactions, transaction)
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
