package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thetsGit/spend-wise-be/internal/presets"
)

func (email *RawEmail) Validate() bool {
	return email.Sender != "" && email.Recipient != "" && email.Subject != "" && email.Date != ""
}

func (s *Spending) Validate() bool {
	return s.Merchant != "" && s.Amount != nil
}

func (s *SaaSDiscovery) Validate() bool {
	return s.ProductName != ""
}

func (spending *AISpendingResult) CalculateScore() string {
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

func (saas *AISaaSDiscoveryResult) CalculateScore() string {
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

// Parse raw AI response body, removing trailing prefix/postfix if exists
func ParseAIResponse(raw string) ([]AIResult, error) {

	parsed := strings.TrimSpace(raw)
	parsed = strings.TrimPrefix(parsed, "```json")
	parsed = strings.TrimPrefix(parsed, "```")
	parsed = strings.TrimSuffix(parsed, "```")
	parsed = strings.TrimSpace(parsed)

	var results []AIResult
	if err := json.Unmarshal([]byte(parsed), &results); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}
	return results, nil
}
