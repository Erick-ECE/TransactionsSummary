package usecases

import (
	"fmt"
	"log"
	"strconv"
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
func (uc *SendSummaryEmail) Execute(accountToTransactions map[string][]entities.Transaction) error {
	// Generate the summary
	for account, transactions := range accountToTransactions {
		summaryResult, toEmail, err := uc.GenerateSummaryUseCase.Execute(account, transactions)

		if err != nil {
			log.Printf("Could not generate summary for account %s: %v", account, err)
			return fmt.Errorf("could not generate summary: %v", err)
		}

		// Format the summary into HTML
		emailBody := uc.formatSummaryAsHTML(summaryResult)

		// Send the email
		subject := "Monthly Transactions Summary"
		if err := uc.EmailSender.SendEmail(toEmail, subject, emailBody); err != nil {
			log.Printf("Could not send summary email to %s: %v", toEmail, err)
			return fmt.Errorf("could not send summary email: %v", err)
		}
	}

	return nil
}

// formatSummaryAsHTML converts the summary result to an HTML format.
func (uc *SendSummaryEmail) formatSummaryAsHTML(summary *entities.SummaryResult) string {
	var sb strings.Builder

	// Begin email container
	sb.WriteString(`
    <div style="background-color: #ffffff; max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif;">
        <!-- Header with logo -->
        <div style="background-color: #b9ff66; text-align: center; padding: 20px;">
            <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/b/b0/Stori_Logo_2023.svg/512px-Stori_Logo_2023.svg.png" 
                 alt="Stori Logo" 
                 style="width: 150px; height: auto;">
        </div>

        <!-- Main Content -->
        <div style="padding: 30px 40px;">
            <h1 style="color: #000000; font-size: 24px; margin-bottom: 20px;">Transactions Summary</h1>
            
            <!-- Total Summary Cards -->
            <div style="display: inline-block; width: 45%; margin-right: 5%; background-color: #f8f9fa; padding: 15px; border-radius: 8px;">
                <h3 style="margin: 0; color: #666;">Total Credit</h3>
                <p style="font-size: 24px; margin: 10px 0; color: #28a745;">`)
	sb.WriteString(strconv.FormatFloat(summary.TotalCredit, 'f', 2, 64))
	sb.WriteString(`</p>
            </div>
            <div style="display: inline-block; width: 45%; background-color: #f8f9fa; padding: 15px; border-radius: 8px;">
                <h3 style="margin: 0; color: #666;">Total Debit</h3>
                <p style="font-size: 24px; margin: 10px 0; color: #dc3545;">`)
	sb.WriteString(strconv.FormatFloat(summary.TotalDebit, 'f', 2, 64))
	sb.WriteString(`</div>

            <!-- Monthly Breakdown -->
            <h2 style="color: #000000; font-size: 20px; margin: 30px 0 20px;">Monthly Breakdown</h2>
            <table style="width: 100%; border-collapse: collapse; margin-bottom: 30px;">
                <thead>
                    <tr style="background-color: #b9ff66;">
                        <th style="padding: 12px; text-align: left; border-bottom: 2px solid #dee2e6;">Month</th>
                        <th style="padding: 12px; text-align: right; border-bottom: 2px solid #dee2e6;">Transactions</th>
                        <th style="padding: 12px; text-align: right; border-bottom: 2px solid #dee2e6;">Avg Credit</th>
                        <th style="padding: 12px; text-align: right; border-bottom: 2px solid #dee2e6;">Avg Debit</th>
                    </tr>
                </thead>
                <tbody>`)

	// Add monthly rows
	for _, monthSummary := range summary.MonthlySummaries {
		sb.WriteString(fmt.Sprintf(`
                    <tr style="border-bottom: 1px solid #dee2e6;">
                        <td style="padding: 12px; text-align: left;">%s</td>
                        <td style="padding: 12px; text-align: center;">%d</td>
                        <td style="padding: 12px; text-align: right;">$%.2f</td>
                        <td style="padding: 12px; text-align: right;">$%.2f</td>
                    </tr>`,
			monthSummary.Month,
			monthSummary.NumTransactions,
			monthSummary.AverageCredit,
			monthSummary.AverageDebit))
	}

	// Close the table and add footer
	sb.WriteString(`
                </tbody>
            </table>
        </div>

        <!-- Footer -->
        <div style="background-color: #f8f9fa; padding: 20px; text-align: center;">
        
            <p style="color: #666; font-size: 12px; margin: 0;">
                © 2024 Stori. All rights reserved.<br>
                <a href="#" style="color: #666; text-decoration: none;">Privacy Policy</a> | 
                <a href="#" style="color: #666; text-decoration: none;">Unsubscribe</a>
            </p>
        </div>
    </div>`)

	return sb.String()
}
