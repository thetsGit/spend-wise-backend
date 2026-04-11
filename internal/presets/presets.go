package presets

var SpendingCategories = map[string]string{
	"food_delivery": "Food ordering and delivery (GrabFood, foodpanda, Uber Eats)",
	"travel":        "Transportation like taxis, ride-hailing (Uber, Grab), flights, hotels",
	"software":      "SaaS subscriptions, software licenses, digital tools (e.g Zoom, Microsoft Office, etc)",
	"shopping":      "Physical product purchases, e-commerce orders",
	"utilities":     "Electricity, water, internet, phone bills",
	"entertainment": "Streaming services (Netflix, Spotify), games, media, etc",
	"other":         "Anything that does not fit above",
}

var SaaSSignalTypes = map[string]string{
	"welcome":        "Onboarding or signup confirmation email from a SaaS product",
	"invoice":        "Payment receipt or billing confirmation for a subscription",
	"renewal":        "Notification that a subscription has been automatically renewed",
	"trial_expiring": "Warning that a free trial is ending soon, usually with upgrade prompt",
	"usage_report":   "Periodic summary of product usage, activity, or analytics",
}

var SaaSBillingCycles = map[string]string{
	"monthly":   "Billed once per month",
	"quarterly": "Billed once every three months",
	"yearly":    "Billed once per year, sometimes shown as annual",
	"unknown":   "Billing frequency cannot be determined from the email",
}

const (
	ConfidenceHigh   = "high"
	ConfidenceMedium = "medium"
	ConfidenceLow    = "low"
)

var ConfidenceScores = map[string]string{
	ConfidenceHigh:   "All key fields extracted clearly from the email content",
	ConfidenceMedium: "Some fields extracted but one or more are uncertain or missing",
	ConfidenceLow:    "Very little could be extracted, most fields are null or guessed",
}

const (
	EmailStatusPending   = "pending"
	EmailStatusProcessed = "processed"
	EmailStatusFailed    = "failed"
)

var EmailStatuses = map[string]string{
	EmailStatusPending:   "Email uploaded but not yet processed by AI",
	EmailStatusProcessed: "AI has analyzed the email and results are stored",
	EmailStatusFailed:    "AI processing failed for this email",
}
