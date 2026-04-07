package models

import "time"

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

type AISaaSDiscoveryResult struct {
	ProductName   string   `json:"product_name"`
	SignalType    string   `json:"signal_type"`
	BillingCycle  string   `json:"billing_cycle"`
	EstimatedCost *float64 `json:"estimated_cost"`
	Currency      string   `json:"currency"`
	Confidence    string   `json:"confidence"`
}

type SaaSDiscovery struct {
	ID            int
	ProductName   string
	SignalType    string
	BillingCycle  string
	EstimatedCost *float64
	Currency      string
	AIConfidence  *string
	Confidence    string
	CreatedAt     time.Time
	EmailID       int
}

type AISaaSSpendingResult struct {
	Merchant        string    `json:"merchant"`
	Amount          *float64  `json:"amount"`
	Currency        string    `json:"currency"`
	Category        string    `json:"category"`
	TransactionDate time.Time `json:"transaction_date"`
	Confidence      *string   `json:"confidence"`
}

type Spending struct {
	ID              int
	EmailID         int
	Merchant        string
	Amount          *float64
	Currency        string
	Category        string
	TransactionDate time.Time
	AIConfidence    *string
	Confidence      string
	CreatedAt       time.Time
}
