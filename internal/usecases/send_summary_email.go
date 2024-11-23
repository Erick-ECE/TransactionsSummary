package usecases

import (
	"fmt"
	"strings"

	"transactions-summary/internal/entities"
	"transactions-summary/internal/interfaces"
)

// SendSummaryEmail is a use case that generates a summary and sends it via email.
type SendSummaryEmail struct {
	GenerateSummaryUseCase *GenerateSummary
	EmailSender            interfaces.EmailSender
}

// NewSendSummaryEmail creates a new SendSummaryEmail use case.
func NewSendSummaryEmail(generateSummary *GenerateSummary, emailSender interfaces.EmailSender) *SendSummaryEmail {
	return &SendSummaryEmail{
		GenerateSummaryUseCase: generateSummary,
		EmailSender:            emailSender,
	}
}

// Execute generates the summary and sends it to the specified email address.
func (uc *SendSummaryEmail) Execute(toEmail string) error {
	// Generate the summary
	summaryResult, err := uc.GenerateSummaryUseCase.Execute()
	if err != nil {
		return fmt.Errorf("could not generate summary: %v", err)
	}

	// Format the summary into HTML
	emailBody := uc.formatSummaryAsHTML(summaryResult)

	// Send the email
	subject := "Monthly Transactions Summary"
	if err := uc.EmailSender.SendEmail(toEmail, subject, emailBody); err != nil {
		return fmt.Errorf("could not send summary email: %v", err)
	}

	return nil
}

// formatSummaryAsHTML converts the summary result to an HTML format.
func (uc *SendSummaryEmail) formatSummaryAsHTML(summary *entities.SummaryResult) string {
	var sb strings.Builder
	sb.WriteString("<h1>Monthly Transactions Summary</h1>")
	sb.WriteString(fmt.Sprintf("<p>Total Balance: <strong>%.2f</strong></p>", summary.TotalBalance))
	sb.WriteString("<h2>Monthly Breakdown:</h2>")

	sb.WriteString("<table border='1' cellpadding='5' cellspacing='0'>")
	sb.WriteString("<tr><th>Month</th><th>Transactions</th><th>Average Credit</th><th>Average Debit</th></tr>")
	for _, monthSummary := range summary.MonthlySummaries {
		sb.WriteString("<tr>")
		sb.WriteString(fmt.Sprintf("<td>%s</td>", monthSummary.Month))
		sb.WriteString(fmt.Sprintf("<td>%d</td>", monthSummary.NumTransactions))
		sb.WriteString(fmt.Sprintf("<td>%.2f</td>", monthSummary.AverageCredit))
		sb.WriteString(fmt.Sprintf("<td>%.2f</td>", monthSummary.AverageDebit))
		sb.WriteString("</tr>")
	}
	sb.WriteString("</table>")

	return sb.String()
}
