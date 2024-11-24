package entities

// MonthlySummary holds the summary information for a single month.
type MonthlySummary struct {
	Month           string  // E.g., "July"
	NumTransactions int     // Number of transactions in the month
	AverageCredit   float64 // Average credit amount
	AverageDebit    float64 // Average debit amount
	TotalCredits    float64 // Total of all credit transactions
	TotalDebits     float64 // Total of all debit transactions
}

// SummaryResult holds the overall summary data.
type SummaryResult struct {
	TotalCredit      float64
	TotalDebit       float64
	MonthlySummaries []MonthlySummary // Summary grouped by month

}
