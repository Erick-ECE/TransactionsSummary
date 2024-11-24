package entities

import "time"

// represents a transaction (debit or credit)
type Transaction struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	Type            string    `json:"type"` // "debit" or "credit"
}
