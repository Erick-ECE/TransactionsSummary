package database

import (
	"database/sql"
	"fmt"
	"time"

	"transactions-summary/internal/entities"
	"transactions-summary/internal/interfaces"

	_ "github.com/go-sql-driver/mysql" // MySQL driver for database connection
)

// MySQLTransactionRepo implements the TransactionRepository interface for MySQL.
type MySQLTransactionRepo struct {
	DB *sql.DB
}

// Ensure MySQLTransactionRepo implements interfaces.TransactionRepository
var _ interfaces.TransactionRepository = &MySQLTransactionRepo{}

// NewMySQLTransactionRepo creates a new MySQLTransactionRepo instance.
func NewMySQLTransactionRepo(db *sql.DB) *MySQLTransactionRepo {
	return &MySQLTransactionRepo{DB: db}
}

// SaveTransaction saves a new transaction to the database.
func (repo *MySQLTransactionRepo) SaveTransaction(transaction *entities.Transaction) error {
	_, err := repo.DB.Exec(
		"INSERT INTO transactions (id, amount, date, type) VALUES (?, ?, ?, ?)",
		transaction.ID, transaction.Amount, transaction.Date, transaction.Type,
	)
	if err != nil {
		return fmt.Errorf("could not save transaction: %v", err)
	}
	return nil
}

// GetAllTransactions retrieves all transactions from the database.
func (repo *MySQLTransactionRepo) GetAllTransactions() ([]entities.Transaction, error) {
	rows, err := repo.DB.Query("SELECT id, amount, date, type FROM transactions")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve transactions: %v", err)
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction
		var dateString string // Scan the date as a string instead of time.Time
		err := rows.Scan(&transaction.ID, &transaction.Amount, &dateString, &transaction.Type)
		if err != nil {
			return nil, fmt.Errorf("could not scan transaction row: %v", err)
		}

		// Convert the dateString to time.Time
		transaction.Date, err = time.Parse("2006-01-02", dateString)
		if err != nil {
			return nil, fmt.Errorf("could not parse date: %v", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
