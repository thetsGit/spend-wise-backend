package presets

var SpendingCategories = map[string]bool{
	"food_delivery": true,
	"travel":        true,
	"software":      true,
	"shopping":      true,
	"utilities":     true,
	"entertainment": true,
	"other":         true,
}

var SaaSSignalTypes = map[string]bool{
	"welcome":        true,
	"invoice":        true,
	"renewal":        true,
	"trial_expiring": true,
	"usage_report":   true,
}

var SaaSBillingCycles = map[string]bool{
	"monthly":   true,
	"quarterly": true,
	"yearly":    true,
	"unknown":   true,
}

const (
	ConfidenceHigh   = "high"
	ConfidenceMedium = "medium"
	ConfidenceLow    = "low"
)

var ConfidenceScores = map[string]bool{
	ConfidenceHigh:   true,
	ConfidenceMedium: true,
	ConfidenceLow:    true,
}

const (
	EmailStatusPending   = "pending"
	EmailStatusProcessed = "processed"
	EmailStatusFailed    = "failed"
)

var EmailStatuses = map[string]bool{
	EmailStatusPending:   true,
	EmailStatusProcessed: true,
	EmailStatusFailed:    true,
}
