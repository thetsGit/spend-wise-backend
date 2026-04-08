package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/ai"
	"github.com/thetsGit/spend-wise-be/internal/models"
	"github.com/thetsGit/spend-wise-be/internal/presets"
	"github.com/thetsGit/spend-wise-be/internal/prompts"
)

func (h *Handler) UploadEmails(w http.ResponseWriter, r *http.Request) {

	/**
	 * This handler triggers email expense analyzer pipeline step by step
	 */

	/**
	 * (1) Parse JSON body
	 */

	/**
	 * (1.1) Guard against large input sizes
	 */

	// Check content length first
	if r.ContentLength > h.Config.MaxUploadSizeBytes {
		RespondErrorJSON(w, "File too large", http.StatusRequestEntityTooLarge, nil)
		return
	}

	// Then, check body size as well
	r.Body = http.MaxBytesReader(w, r.Body, h.Config.MaxUploadSizeBytes)

	json.Marshal(r.Body)

	var emails []models.RawEmail
	err := json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		RespondErrorJSON(w, "Request body too large or invalid JSON", http.StatusBadRequest, err)
		return
	}

	/**
	 * (2) Validate the emails
	 */

	var validEmails []models.RawEmail
	for _, email := range emails {
		if email.Validate() {
			validEmails = append(validEmails, email)
		}
	}

	/**
	 * (3) Save validate emails to database
	 */

	var formattedEmails []models.Email
	for _, email := range validEmails {
		// TODO: should add batch transaction controll to rollback ?
		email, err := h.DB.InsertEmail(email)
		if err != nil {
			RespondErrorJSON(w, "Failed to save emails", http.StatusInternalServerError, err)
			return
		}
		formattedEmails = append(formattedEmails, email)
	}

	/**
	 * (4) Build AI prompt for spending, saas records retrieval
	 */

	prompt := prompts.BuildPrompt(formattedEmails)

	/**
	 * (5) Trigger AI call with the prepared prompt
	 */

	rawAIResult, err := ai.CallOpenAI(prompt, h.Config)

	if err != nil {
		RespondErrorJSON(w, "AI Failed to process the emails", http.StatusInternalServerError, err)
		return
	}

	/**
	 * (6) Parse the AI responses to retrive spending, saas records
	 */

	aiResult, err := models.ParseAIResponse(rawAIResult)
	if err != nil {
		RespondErrorJSON(w, "Failed to parse AI results", http.StatusInternalServerError, err)
		return
	}

	/**
	 * (7) Normalize, calculate scores on spending, saas records, and save them to database
	 */

	var savedSpendingCount int
	var savedSaaSDiscoveryCount int

	for _, result := range aiResult {
		_, err := h.DB.UpdateEmailStatus(result.EmailID, presets.EmailStatusProcessed)

		if err != nil {
			RespondErrorJSON(w, "Failed to update email status", http.StatusInternalServerError, err)
			return
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

			_, err := h.DB.InsertSpending(spending)

			if err != nil {
				RespondErrorJSON(w, "Failed to save spending record", http.StatusInternalServerError, err)
				return
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

			_, err := h.DB.InsertSaaSDiscovery(saasDiscovery)

			if err != nil {
				RespondErrorJSON(w, "Failed to save saas discovery record", http.StatusInternalServerError, err)
				return
			}

			savedSaaSDiscoveryCount++
		}

	}

	/**
	 * (8) Evaluate summary data
	 */

	var summary models.UploadSummary

	summary.TotalEmails = len(emails)
	summary.Inserted = len(validEmails)
	summary.Skipped = summary.TotalEmails - summary.Inserted

	summary.SpendingFound = savedSpendingCount
	summary.SaaSFound = savedSaaSDiscoveryCount

	/**
	 * (9) return the summary data and end the pipeline
	 */

	RespondDataJSON(w, "Emails processed", http.StatusOK, summary)
}
