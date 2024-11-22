package usecases

import (
	"fmt"

	"transactions-summary/internal/entities"
	"transactions-summary/internal/interfaces"
)

// GenerateSummary processes transactions and creates a summary.
type GenerateSummary struct {
	TransactionRepo interfaces.TransactionRepository
}

// NewGenerateSummary creates a new GenerateSummary use case.
func NewGenerateSummary(repo interfaces.TransactionRepository) *GenerateSummary {
	return &GenerateSummary{
		TransactionRepo: repo,
	}
}

// Execute calculates the summary from all transactions in the database.
func (uc *GenerateSummary) Execute() (*entities.SummaryResult, error) {
	// Retrieve all transactions from the repository
	transactions, err := uc.TransactionRepo.GetAllTransactions()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve transactions: %v", err)
	}

	// Calculate summary data
	totalBalance := 0.0
	monthlyData := make(map[string]*entities.MonthlySummary)

	// Process each transaction
	for _, transaction := range transactions {
		// Update the total balance
		totalBalance += transaction.Amount

		// Get month name (e.g., "July")
		monthName := transaction.Date.Format("January")

		// Initialize monthly summary if not present
		if _, exists := monthlyData[monthName]; !exists {
			monthlyData[monthName] = &entities.MonthlySummary{
				Month: monthName,
			}
		}

		// Update monthly data
		monthlySummary := monthlyData[monthName]
		monthlySummary.NumTransactions++

		if transaction.Type == "credit" {
			monthlySummary.TotalCredits += transaction.Amount
		} else if transaction.Type == "debit" {
			monthlySummary.TotalDebits += transaction.Amount
		}
	}

	// Calculate averages for each month
	var monthlySummaries []entities.MonthlySummary
	for _, summary := range monthlyData {
		if summary.NumTransactions > 0 {
			// Calculate averages
			creditCount := float64(countCredits(transactions, summary.Month))
			debitCount := float64(countDebits(transactions, summary.Month))

			if creditCount > 0 {
				summary.AverageCredit = summary.TotalCredits / creditCount
			}
			if debitCount > 0 {
				summary.AverageDebit = summary.TotalDebits / debitCount
			}
		}
		monthlySummaries = append(monthlySummaries, *summary)
	}

	return &entities.SummaryResult{
		TotalBalance:     totalBalance,
		MonthlySummaries: monthlySummaries,
	}, nil
}

// Helper functions for counting transactions by type
func countCredits(transactions []entities.Transaction, month string) int {
	count := 0
	for _, transaction := range transactions {
		if transaction.Date.Format("January") == month && transaction.Type == "credit" {
			count++
		}
	}
	return count
}

func countDebits(transactions []entities.Transaction, month string) int {
	count := 0
	for _, transaction := range transactions {
		if transaction.Date.Format("January") == month && transaction.Type == "debit" {
			count++
		}
	}
	return count
}
