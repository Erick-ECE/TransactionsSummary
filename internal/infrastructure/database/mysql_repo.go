package database

import (
	"database/sql"
	"fmt"
	"log"
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
func (repo *MySQLTransactionRepo) SaveTransaction(transaction entities.Transaction) error {
	_, err := repo.DB.Exec(
		"INSERT INTO transactions (id, account_id, amount, transaction_date, type) VALUES (?, ?, ?, ?, ?)",
		transaction.ID, transaction.AccountID, transaction.Amount, transaction.TransactionDate, transaction.Type,
	)
	if err != nil {
		log.Printf("Error saving transaction %s: %v", transaction.ID, err)
		return fmt.Errorf("could not save transaction: %v", err)
	}
	log.Printf("Transaction %s saved successfully", transaction.ID)
	return nil
}

// GetTransaction retrieves a transaction from the database by ID.
func (repo *MySQLTransactionRepo) GetTransaction(transactionID string) (*entities.Transaction, error) {
	query := "SELECT id, account_id, amount, transaction_date, type FROM transactions WHERE id = ?"

	// Create a variable to hold the account details
	transaction := &entities.Transaction{}
	var dateString string

	// Execute the query and scan the result into the account struct
	err := repo.DB.QueryRow(query, transactionID).Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &dateString, &transaction.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Transaction with ID %s not found", transactionID)
			return nil, fmt.Errorf("transaction with id %s not found", transactionID)
		}
		log.Printf("Error retrieving transaction %s: %v", transactionID, err)
		return nil, fmt.Errorf("could not retrieve account: %v", err)
	}

	// Convert the dateString to time.Time
	transaction.TransactionDate, err = time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}

	return transaction, nil
}

// GetAccount retrieves an account from the database by ID.
func (repo *MySQLTransactionRepo) GetAccount(id string) (*entities.Account, error) {
	query := "SELECT id, debit_balance, credit_balance, email FROM accounts WHERE id = ?"

	// Create a variable to hold the account details
	account := &entities.Account{}

	// Execute the query and scan the result into the account struct
	err := repo.DB.QueryRow(query, id).Scan(&account.ID, &account.DebitBalance, &account.CreditBalance, &account.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Account with ID %s not found", id)
			return nil, fmt.Errorf("account with id %s not found", id)
		}
		log.Printf("Error retrieving account %s: %v", id, err)
		return nil, fmt.Errorf("could not retrieve account: %v", err)
	}

	log.Printf("Account %s retrieved successfully", id)
	return account, nil
}

// UpdateAccount updates a given account from the database.
func (repo *MySQLTransactionRepo) UpdateAccount(account *entities.Account) error {
	_, err := repo.DB.Exec(
		"UPDATE accounts SET debit_balance = ?, credit_balance = ? WHERE id = ?", account.DebitBalance, account.CreditBalance, account.ID,
	)
	if err != nil {
		log.Printf("Error updating account %s: %v", account.ID, err)
		return fmt.Errorf("could not save transaction: %v", err)
	}
	log.Printf("Account %s updated successfully", account.ID)
	return nil
}
