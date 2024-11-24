package entities

type Account struct {
	ID            string  `json:"id"`
	DebitBalance  float64 `json:"debit_balance"`
	CreditBalance float64 `json:"credit_balance"`
	Email         string  `json:"email"`
}
