package interfaces

import "transactions-summary/internal/entities"

// FileReader defines the interface for reading transactions from a file.
type FileReader interface {
	ReadTransactions(filePath string) ([]entities.Transaction, error)
}
