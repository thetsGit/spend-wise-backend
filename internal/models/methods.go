package models

import "github.com/thetsGit/spend-wise-be/internal/presets"

func (email *Email) Validate() bool {
	return email.Sender != "" && email.Recipient != "" && email.Subject != "" && !email.Date.IsZero()
}

func (spending *Spending) CalculateScore() string {
	const MAX_SCORE = 5
	score := 0

	if spending.Merchant != "" {
		score++
	}

	if spending.Amount != nil {
		score++
	}

	if spending.Currency != "" {
		score++
	}

	if spending.Category != "" {
		score++
	}

	if spending.TransactionDate != nil {
		score++
	}

	if score == MAX_SCORE {
		return presets.ConfidenceHigh
	}

	if score == 0 {
		return presets.ConfidenceLow
	}

	return presets.ConfidenceMedium
}

func (saas *SaaSDiscovery) CalculateScore() string {
	const MAX_SCORE = 5
	score := 0

	if saas.ProductName != "" {
		score++
	}

	if saas.SignalType != "" {
		score++
	}

	if saas.BillingCycle != "" {
		score++
	}

	if saas.EstimatedCost != nil {
		score++
	}

	if saas.Currency != "" {
		score++
	}

	if score == MAX_SCORE {
		return presets.ConfidenceHigh
	}

	if score == 0 {
		return presets.ConfidenceLow
	}

	return presets.ConfidenceMedium
}
