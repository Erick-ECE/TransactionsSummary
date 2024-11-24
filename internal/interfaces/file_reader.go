package interfaces

import (
	"encoding/csv"

	"transactions-summary/internal/entities"
)

// FileReader defines the interface for reading transactions from a file.
type FileReader interface {
	ReadTransactions(reader *csv.Reader) ([]entities.Transaction, error)
}
