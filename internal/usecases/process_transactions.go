package usecases

import (
	"encoding/csv"
	"fmt"

	"transactions-summary/internal/entities"
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
func (uc *ProcessTransactions) Execute(reader *csv.Reader) (map[string][]entities.Transaction, error) {
	// Read the transactions from the file
	transactions, err := uc.FileReader.ReadTransactions(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read transactions: %v", err)
	}

	var filteredTransaction []entities.Transaction

	for _, transaction := range transactions {
		t11n, _ := uc.TransactionRepo.GetTransaction(transaction.ID)
		if t11n != nil {
			continue
		}
		filteredTransaction = append(filteredTransaction, transaction)
	}

	// Process each transaction and save to the database
	for _, txn := range filteredTransaction {

		// Save the transaction to the database
		err = uc.TransactionRepo.SaveTransaction(txn)
		if err != nil {
			return nil, fmt.Errorf("could not save transaction: %v", err)
		}
	}

	accountsToTransaction := make(map[string][]entities.Transaction)

	for _, transaction := range filteredTransaction {
		accountsToTransaction[transaction.AccountID] = append(accountsToTransaction[transaction.AccountID], transaction)
	}

	return accountsToTransaction, nil
}
