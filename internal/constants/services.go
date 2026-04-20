package constants

import "time"

const (
	// Open AI services
	OpenAIUrl   = "https://api.openai.com/v1/chat/completions"
	OpenAIModel = "gpt-5.4"

	// Google services
	GoogleOAuthURL       = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL    = "https://www.googleapis.com/oauth2/v2/userinfo"
	GmailAPIURL          = "https://gmail.googleapis.com/gmail/v1"
	OpenAIURL            = "https://api.anthropic.com/v1/messages"
	MaxGmailCountPerSync = "30"

	// Auth services
	AuthSessionLifeSec = 7 * 24 * time.Hour // 7 days

)
