package usecases

import (
	"testing"
	"time"

	"transactions-summary/internal/entities"

	"github.com/stretchr/testify/assert"
)

// MockTransactionRepo is a mock implementation of the TransactionRepository interface.
type MockTransactionRepo struct {
	transactions []entities.Transaction
}

// NewMockTransactionRepo creates a new mock repository with predefined transactions.
func NewMockTransactionRepo(transactions []entities.Transaction) *MockTransactionRepo {
	return &MockTransactionRepo{transactions: transactions}
}

// GetAllTransactions returns the predefined list of transactions.
func (m *MockTransactionRepo) GetAllTransactions() ([]entities.Transaction, error) {
	return m.transactions, nil
}

// SaveTransaction is a mock method for saving a transaction (not used in tests).
func (m *MockTransactionRepo) SaveTransaction(transaction *entities.Transaction) error {
	return nil
}

func TestGenerateSummary_Execute(t *testing.T) {
	// Define test data: transactions for different months with credits and debits
	testTransactions := []entities.Transaction{
		{ID: "1", Amount: 60.5, Date: parseDate("2024-07-15"), Type: "credit"},
		{ID: "2", Amount: -10.3, Date: parseDate("2024-07-28"), Type: "debit"},
		{ID: "3", Amount: -20.46, Date: parseDate("2024-08-02"), Type: "debit"},
		{ID: "4", Amount: 10.0, Date: parseDate("2024-08-13"), Type: "credit"},
		{ID: "5", Amount: 27.0, Date: parseDate("2024-10-15"), Type: "credit"},
	}

	// Initialize the mock repository with test transactions
	mockRepo := NewMockTransactionRepo(testTransactions)

	// Initialize the GenerateSummary use case
	summaryUseCase := NewGenerateSummary(mockRepo)

	// Execute the summary calculation
	summaryResult, err := summaryUseCase.Execute()
	assert.NoError(t, err, "expected no error during summary generation")

	// Validate the results
	assert.NotNil(t, summaryResult)
	assert.Equal(t, 66.74, summaryResult.TotalBalance) // Validate the total balance

	// Validate monthly summaries
	assert.Len(t, summaryResult.MonthlySummaries, 3) // Should have 3 unique months

	// July data validation
	julySummary := findMonthlySummary(summaryResult.MonthlySummaries, "July")
	assert.Equal(t, 2, julySummary.NumTransactions)
	assert.Equal(t, 60.5, julySummary.AverageCredit)
	assert.Equal(t, -10.3, julySummary.AverageDebit)

	// August data validation
	augustSummary := findMonthlySummary(summaryResult.MonthlySummaries, "August")
	assert.Equal(t, 2, augustSummary.NumTransactions)
	assert.Equal(t, 10.0, augustSummary.AverageCredit)
	assert.Equal(t, -20.46, augustSummary.AverageDebit)

	// October data validation
	octoberSummary := findMonthlySummary(summaryResult.MonthlySummaries, "October")
	assert.Equal(t, 1, octoberSummary.NumTransactions)
	assert.Equal(t, 27.0, octoberSummary.AverageCredit)
	assert.Equal(t, 0.0, octoberSummary.AverageDebit) // No debits in October
}

// parseDate is a helper function to parse a date string into a time.Time object.
func parseDate(dateStr string) time.Time {
	date, _ := time.Parse("2006-01-02", dateStr)
	return date
}

// findMonthlySummary is a helper function to locate a MonthlySummary by month name.
func findMonthlySummary(summaries []entities.MonthlySummary, month string) *entities.MonthlySummary {
	for _, summary := range summaries {
		if summary.Month == month {
			return &summary
		}
	}
	return nil
}
