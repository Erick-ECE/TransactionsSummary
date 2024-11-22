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
		var txn entities.Transaction
		var date time.Time
		err := rows.Scan(&txn.ID, &txn.Amount, &date, &txn.Type)
		if err != nil {
			return nil, fmt.Errorf("could not scan transaction row: %v", err)
		}
		txn.Date = date
		transactions = append(transactions, txn)
	}

	return transactions, nil
}
