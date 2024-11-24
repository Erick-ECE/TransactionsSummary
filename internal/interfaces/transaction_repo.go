package interfaces

import "transactions-summary/internal/entities"

// TransactionRepository defines the interface for database operations.
type TransactionRepository interface {
	SaveTransaction(transaction entities.Transaction) error
	GetAccount(accountId string) (*entities.Account, error)
	UpdateAccount(account *entities.Account) error
	GetTransaction(transactionID string) (*entities.Transaction, error)
}
