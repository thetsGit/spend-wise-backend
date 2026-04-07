package prompts

import (
	"fmt"
	"strings"

	"github.com/thetsGit/spend-wise-be/internal/models"
	"github.com/thetsGit/spend-wise-be/internal/presets"
	"github.com/thetsGit/spend-wise-be/internal/utils"
)

func BuildPrompt(emails []models.Email) string {
	categories := strings.Join(utils.Keys(presets.SpendingCategories), ", ")
	signalTypes := strings.Join(utils.Keys(presets.SaaSSignalTypes), ", ")
	billingCycles := strings.Join(utils.Keys(presets.SaaSBillingCycles), ", ")
	confidences := strings.Join(utils.Keys(presets.ConfidenceScores), ", ")

	prompt := fmt.Sprintf(`You are an email analyzer that performs two tasks:
1. Extract financial transactions (spending)
2. Detect SaaS/software product signals

Rules for spending:
- category must be one of: %s
- amount should be a number without currency symbols
- date should be in ISO 8601 format. Extract the transaction date from email content, not the email date
- currency should be a 3-letter code. Default to USD if unclear

Rules for SaaS:
- signal_type must be one of: %s
- billing_cycle must be one of: %s
- estimated_cost should be the total amount (not per-user), as a number without currency symbols
- currency should be a 3-letter code. Default to USD if unclear

General rules:
- confidence must be one of: %s
- If an email has no spending data, set spending to null
- If an email has no SaaS signal, set saas to null
- If you cannot determine a field, use null

Emails:
`, categories, signalTypes, billingCycles, confidences)

	for i, e := range emails {
		prompt += fmt.Sprintf(`
Email %d:
From: %s
Subject: %s
Date: %s
Body: %s
`, i, e.Sender, e.Subject, e.Date.Format("2006-01-02"), e.Body)
	}

	prompt += `
Respond with ONLY a valid JSON array. No markdown, no explanation.
[
  {
    "email_id": 0,
    "spending": {
      "merchant": "string",
      "amount": 0.00,
      "currency": "USD",
      "category": "string",
      "date": "2025-07-01T00:00:00Z",
      "confidence": "high"
    },
    "saas": {
      "product_name": "string",
      "signal_type": "string",
      "billing_cycle": "string",
      "estimated_cost": 0.00,
      "currency": "USD",
      "confidence": "high"
    }
  }
]`

	return prompt
}
