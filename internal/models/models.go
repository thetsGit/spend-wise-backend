package models

import "time"

/**
 * Business oriented entities (structs)
 */

type RawEmail struct {
	Sender    string `json:"from"`
	Recipient string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Date      string `json:"date"`
}

type Email struct {
	ID        int
	Sender    string
	Recipient string
	Subject   string
	Body      string
	Date      time.Time
	Status    string
	CreatedAt time.Time
}

type AISpendingResult struct {
	Merchant        string     `json:"merchant"`
	Amount          *float64   `json:"amount"`
	Currency        string     `json:"currency"`
	Category        string     `json:"category"`
	TransactionDate *time.Time `json:"date"`
	Confidence      *string    `json:"confidence"`
}

type Spending struct {
	ID              int        `db:"id" json:"id"`
	EmailID         int        `db:"email_id" json:"email_id"`
	Merchant        string     `db:"merchant" json:"merchant"`
	Amount          *float64   `db:"amount" json:"amount"`
	Currency        string     `db:"currency" json:"currency"`
	Category        string     `db:"category" json:"category"`
	TransactionDate *time.Time `db:"transaction_date" json:"transaction_date"`
	AIConfidence    *string    `db:"ai_confidence" json:"ai_confidence"`
	Confidence      string     `db:"confidence" json:"confidence"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

type SpendingFilter struct {
	Category  string
	StartDate string
	EndDate   string
}

type SpendingSummary struct {
	TotalSpend float64           `json:"total_spend"`
	TotalCount int               `json:"total_count"`
	ByCategory []CategorySummary `json:"by_category"`
}

type CategorySummary struct {
	Category   string  `db:"category" json:"category"`
	TotalSpend float64 `db:"total_spend" json:"total_spend"`
	TotalCount int     `db:"total_count" json:"total_count"`
}

type AISaaSDiscoveryResult struct {
	ProductName   string   `json:"product_name"`
	SignalType    string   `json:"signal_type"`
	BillingCycle  string   `json:"billing_cycle"`
	EstimatedCost *float64 `json:"estimated_cost"`
	Currency      string   `json:"currency"`
	Confidence    *string  `json:"confidence"`
}

type SaaSDiscovery struct {
	ID            int       `db:"id" json:"id"`
	ProductName   string    `db:"product_name" json:"product_name"`
	SignalType    string    `db:"signal_type" json:"signal_type"`
	BillingCycle  string    `db:"billing_cycle" json:"billing_cycle"`
	EstimatedCost *float64  `db:"estimated_cost" json:"estimated_cost"`
	Currency      string    `db:"currency" json:"currency"`
	AIConfidence  *string   `db:"ai_confidence" json:"ai_confidence"`
	Confidence    string    `db:"confidence" json:"confidence"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	EmailID       int       `db:"email_id" json:"email_id"`
}

type SaaSDiscoveryFilter struct {
	ProductName string
	SignalType  string
}

type SaaSSummary struct {
	TotalMonthlySpend float64 `json:"total_monthly_spend"`
	TotalToolsFound   int     `json:"total_tools_found"`
}

type AIResult struct {
	EmailID  int                    `json:"email_id"`
	Spending *AISpendingResult      `json:"spending"`
	SaaS     *AISaaSDiscoveryResult `json:"saas"`
}

type UploadSummary struct {
	TotalEmails   int `json:"total_emails"`
	Inserted      int `json:"inserted"`
	Skipped       int `json:"skipped"`
	SpendingFound int `json:"spending_found"`
	SaaSFound     int `json:"saas_found"`
}

/**
 * General structs
 */

type APIResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Data       any    `json:"data,omitempty"`
	Error      string `json:"error,omitempty"`
}
