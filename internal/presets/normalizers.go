package presets

import "github.com/thetsGit/spend-wise-be/internal/utils"

func NormalizeSpendingCategory(raw string) string {
	return utils.Normalize(raw, SpendingCategories, "other")
}

func NormalizeSaaSSignalType(raw string) string {
	return utils.Normalize(raw, SaaSSignalTypes, "other")
}

func NormalizeBillingCycle(raw string) string {
	return utils.Normalize(raw, SaaSBillingCycles, "unknown")
}

func NormalizeConfidenceScore(raw string) string {
	return utils.Normalize(raw, ConfidenceScores, "low")
}

func NormalizeEmailStatus(raw string) string {
	return utils.Normalize(raw, EmailStatuses, "failed")
}
