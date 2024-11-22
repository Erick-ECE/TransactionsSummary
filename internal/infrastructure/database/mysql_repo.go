package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Required for MySQL database connection

	"transactions-summary/internal/entities"
	"transactions-summary/internal/interfaces" // Check if this is being used correctly
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
