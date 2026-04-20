package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/ai"
	"github.com/thetsGit/spend-wise-be/internal/models"
	"github.com/thetsGit/spend-wise-be/internal/presets"
	"github.com/thetsGit/spend-wise-be/internal/prompts"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondErrorJSON(w http.ResponseWriter, message string, statusCode int, err error) {
	// Server side logging for obervability
	fmt.Printf("[ERROR] %s: %v\n", message, err)

	parsedError := ""
	if err != nil {
		parsedError = err.Error()
	}
	RespondJSON(w, http.StatusOK, models.APIResponse{
		Status:     "error",
		Message:    message,
		StatusCode: statusCode,
		Error:      parsedError,
	})
}

func RespondDataJSON(w http.ResponseWriter, message string, statusCode int, data any) {
	RespondJSON(w, http.StatusOK, models.APIResponse{
		Status:     "success",
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	})
}

func AnalyzeEmails(h *Handler, user models.User, rawEmails []models.RawEmail) (message string, statusCode int, data any, e error) {
	/**
	 * This is the email expense analyzer pipeline step by step
	 */

	/**
	 * (1) Validate the emails
	 */

	var validRawEmails []models.RawEmail
	for _, email := range rawEmails {
		if email.Validate() {
			validRawEmails = append(validRawEmails, email)
		}
	}

	/**
	 * (2) Save validated emails to database
	 */

	var savedEmails []models.Email
	for _, email := range validRawEmails {
		email, err := h.DB.InsertEmail(user.ID, email)
		if err != nil {
			return "Failed to save emails", http.StatusInternalServerError, nil, err
		}

		// Duplicated one, skipped
		isSkipped := email.Sender == "" && email.Recipient == "" && err == nil

		if !isSkipped {
			savedEmails = append(savedEmails, email)
		}
	}

	/**
	 * (2.1, BP-1) If there's no saved emails, break the pipeline right away
	 */

	if len(savedEmails) == 0 {
		/**
		 * Evaluate the summary
		 */
		summary := evaluateSummary(len(rawEmails), len(savedEmails), len(validRawEmails), 0, 0)
		return "Emails processed", http.StatusOK, summary, nil
	}

	/**
	 * (3) Build AI prompts for spending, saas records retrieval
	 */

	systemPrompt := prompts.BuildSystemPrompt()
	userPrompt := prompts.BuildUserPrompt(savedEmails)

	/**
	 * (4) Trigger AI call with the prepared prompt
	 */

	rawAIResult, err := ai.CallOpenAI(systemPrompt, userPrompt, h.Config)

	if err != nil {
		return "AI Failed to process the emails", http.StatusInternalServerError, nil, err
	}

	/**
	 * (5) Parse the AI responses to retrive spending, saas records
	 */

	aiResult, err := models.ParseAIResponse(rawAIResult)
	if err != nil {
		return "Failed to parse AI results", http.StatusInternalServerError, nil, err
	}

	/**
	 * (6) Normalize, calculate scores on spending, saas records, and save them to database
	 */

	var savedSpendingCount int
	var savedSaaSDiscoveryCount int

	for _, result := range aiResult {
		_, err := h.DB.UpdateEmailStatus(result.EmailID, presets.EmailStatusProcessed)

		if err != nil {
			return "Failed to update email status", http.StatusInternalServerError, nil, err
		}

		if result.Spending != nil {
			spending := models.Spending{
				EmailID:         result.EmailID,
				Merchant:        result.Spending.Merchant,
				Amount:          result.Spending.Amount,
				Currency:        result.Spending.Currency,
				Category:        presets.NormalizeSpendingCategory(result.Spending.Category),
				TransactionDate: result.Spending.TransactionDate,
				AIConfidence:    result.Spending.Confidence,
				Confidence:      result.Spending.CalculateScore(),
			}

			_, err := h.DB.InsertSpending(user.ID, spending)

			if err != nil {
				return "Failed to save spending record", http.StatusInternalServerError, nil, err
			}

			savedSpendingCount++
		}

		if result.SaaS != nil {
			saasDiscovery := models.SaaSDiscovery{
				EmailID:       result.EmailID,
				ProductName:   result.SaaS.ProductName,
				SignalType:    presets.NormalizeSaaSSignalType(result.SaaS.SignalType),
				BillingCycle:  presets.NormalizeBillingCycle(result.SaaS.BillingCycle),
				EstimatedCost: result.SaaS.EstimatedCost,
				Currency:      result.SaaS.Currency,
				AIConfidence:  result.SaaS.Confidence,
				Confidence:    result.SaaS.CalculateScore(),
			}

			_, err := h.DB.InsertSaaSDiscovery(user.ID, saasDiscovery)

			if err != nil {
				return "Failed to save saas discovery record", http.StatusInternalServerError, nil, err
			}

			savedSaaSDiscoveryCount++
		}

	}

	/**
	 * (7) Evaluate summary data
	 */

	summary := evaluateSummary(len(rawEmails), len(savedEmails), len(validRawEmails), savedSpendingCount, savedSaaSDiscoveryCount)

	/**
	 * (8) return the summary data and end the pipeline
	 */

	return "Emails processed", http.StatusOK, summary, nil
}

func evaluateSummary(rawECount int, insertedECount int, validECount int, spendingCount int, saasCount int) models.UploadSummary {
	var summary models.UploadSummary

	summary.TotalEmails = rawECount
	summary.Inserted = insertedECount
	summary.Invalid = summary.TotalEmails - validECount
	summary.Skipped = summary.TotalEmails - summary.Inserted - summary.Invalid
	summary.SpendingFound = spendingCount
	summary.SaaSFound = saasCount

	return summary
}
