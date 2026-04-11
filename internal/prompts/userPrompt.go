package prompts

import (
	"fmt"

	"github.com/thetsGit/spend-wise-be/internal/models"
)

func BuildUserPrompt(emails []models.Email) string {
	var prompt string
	prompt = "Analyze the following emails:\n"
	for _, e := range emails {
		prompt += fmt.Sprintf(`
Email: %d
From: %s
Subject: %s
Date: %s
Body: %s
`, e.ID, e.Sender, e.Subject, e.Date.Format("2006-01-02"), e.Body)
	}
	return prompt
}
