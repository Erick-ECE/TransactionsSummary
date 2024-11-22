package entities

import "time"

// represents a transaction (debit or credit)
type Transaction struct {
	ID     string    `json:"id"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
	Type   string    `json:"type"` // "debit" or "credit"
}
