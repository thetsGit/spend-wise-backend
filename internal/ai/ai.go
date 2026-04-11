package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/config"
)

func CallOpenAI(systemPrompt, userPrompt string, config *config.Config) (string, error) {

	/**
	 * Prepare request body
	 */

	reqBody := openAIRequest{
		Model: config.OpenAIModel,
		Messages: []openAIMessage{
			{Role: "developer", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonReqData, jsonReqDataErr := json.Marshal(reqBody)

	if jsonReqDataErr != nil {
		return "failed to parse request body", jsonReqDataErr
	}

	/**
	 * Prepare HTTP request to ai provider
	 */

	req, reqErr := http.NewRequest("POST", config.OpenAIUrl, bytes.NewBuffer(jsonReqData))
	if reqErr != nil {
		return "", fmt.Errorf("failed to create request: %w", reqErr)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.OpenAIApiKey))
	req.Header.Set("Content-Type", "application/json")

	/**
	 * Make HTTP request to open AI
	 */

	client := &http.Client{}
	res, resErr := client.Do(req)
	if resErr != nil {
		return "", fmt.Errorf("API request failed: %w", resErr)
	}
	defer res.Body.Close()

	/**
	 * Parse response
	 */

	resBody, resBodyErr := io.ReadAll(res.Body)
	if resBodyErr != nil {
		return "", fmt.Errorf("failed to read response: %w", resBodyErr)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", res.StatusCode, string(resBody))
	}

	var parsedRes openAIResponse
	parsedResErr := json.Unmarshal(resBody, &parsedRes)

	if parsedResErr != nil {
		return "", fmt.Errorf("failed to parse response: %w", parsedResErr)
	}

	if len(parsedRes.Choices) == 0 {
		return "", fmt.Errorf("empty API response")
	}

	return parsedRes.Choices[0].Message.Content, nil
}
