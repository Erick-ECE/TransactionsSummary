package usecases

import (
	"fmt"

	"transactions-summary/internal/interfaces"
)

// ProcessTransactions processes transactions from a CSV file.
type ProcessTransactions struct {
	TransactionRepo interfaces.TransactionRepository
	FileReader      interfaces.FileReader
}

// NewProcessTransactions creates a new ProcessTransactions use case.
func NewProcessTransactions(repo interfaces.TransactionRepository, reader interfaces.FileReader) *ProcessTransactions {
	return &ProcessTransactions{
		TransactionRepo: repo,
		FileReader:      reader,
	}
}

// Execute reads the CSV file, processes each transaction, and saves them to the database.
func (uc *ProcessTransactions) Execute(filePath string) error {
	// Read the transactions from the file
	transactions, err := uc.FileReader.ReadTransactions(filePath)
	if err != nil {
		return fmt.Errorf("could not read transactions: %v", err)
	}

	// Process each transaction and save to the database
	for _, txn := range transactions {
		// Save the transaction to the database
		err = uc.TransactionRepo.SaveTransaction(&txn)
		if err != nil {
			return fmt.Errorf("could not save transaction: %v", err)
		}
	}

	return nil
}
