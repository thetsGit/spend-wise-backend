package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thetsGit/spend-wise-be/internal/constants"
	"github.com/thetsGit/spend-wise-be/internal/models"
	"golang.org/x/sync/errgroup"
)

type TokenRequestResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GmailsMessage struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
}

type GmailsResponse struct {
	Messages           []GmailsMessage `json:"messages"`
	NextPageToken      string          `json:"nextPageToken"`
	ResultSizeEstimate int             `json:"resultSizeEstimate"`
}

type GmailResponsePayloadHeader struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

type GmailResponsePayload struct {
	Headers []GmailResponsePayloadHeader `json:"headers"`
}

type GmailResponseBody struct {
	Data string `json:"data"`
}

type GmailResponse struct {
	ID           string               `json:"id"`
	ThreadId     string               `json:"threadId"`
	InternalDate string               `json:"internalDate"`
	Payload      GmailResponsePayload `json:"payload"`
	Body         GmailResponseBody    `json:"body"`
}

func getHeader(headers []GmailResponsePayloadHeader, name string) string {
	for _, h := range headers {
		if h.Name == name {
			return h.Value
		}
	}
	return ""
}

func (h *Handler) UploadEmails(w http.ResponseWriter, r *http.Request) {

	/**
	 * Parse JSON body
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

	var rawEmails []models.RawEmail
	err := json.NewDecoder(r.Body).Decode(&rawEmails)
	if err != nil {
		RespondErrorJSON(w, "Request body too large or invalid JSON", http.StatusBadRequest, err)
		return
	}

	/**
	 * Trigger email analyzer pipeline
	 */
	user := GetUserFromContext(r.Context())
	message, statusCode, data, err := AnalyzeEmails(h, *user, rawEmails)

	if err != nil {
		RespondErrorJSON(w, message, statusCode, err)
		return
	}

	RespondDataJSON(w, message, statusCode, data)
}

func (h *Handler) SyncGmails(w http.ResponseWriter, r *http.Request) {
	// fetch latest  x (probably 30) gmails from gmail api
	// use 'access token' stored in db
	// if 'access token' is no longer active, try to refresh it using stored 'refresh token'
	// if failed to refresh, -> respond (error - 401 unauthorized or any relevant code)
	// otherwise (i.e if success),
	// - fetch gmail list
	// - fetch gmail message for each gmail's details
	// - pass to analyzer pipeline

	user := GetUserFromContext(r.Context())
	accessToken := user.OauthAccessToken

	if user.OauthTokenExpiry.Before(time.Now()) {
		// TODO: consider returning statusCode more explicitly (401, 503, etc)

		tokenResponse, refreshErr := h.getNewAccessToken(user.OauthRefreshToken)

		if refreshErr != nil {
			RespondErrorJSON(w, "Failed to refresh token", http.StatusUnauthorized, refreshErr)
			// TODO: backend logout ?
			return
		}

		// Update access token after refresh
		accessToken = tokenResponse.AccessToken

		tokenUpdateErr := h.DB.UpdateAccessToken(user.ID, accessToken, time.Now().Add(time.Duration(tokenResponse.ExpiresIn)*time.Second))

		if tokenUpdateErr != nil {
			RespondErrorJSON(w, "Failed to update access token", http.StatusInternalServerError, tokenUpdateErr)
			// TODO: backend logout ?
			return
		}
	}

	rMessages, getGmailsErr := h.getGmails(accessToken)

	if getGmailsErr != nil {
		RespondErrorJSON(w, "Failed to get gmails", http.StatusServiceUnavailable, getGmailsErr)
		return
	}

	/**
	 * Fetch message details
	 */

	group, ctx := errgroup.WithContext(context.Background())
	messages := make([]*GmailResponse, len(rMessages.Messages))

	for i, message := range rMessages.Messages {
		// Parallel fetch using go routines
		group.Go(func() error {
			msg, err := h.getGmailMessage(message.ID, accessToken, ctx)
			if err != nil {
				return err
			}
			messages[i] = msg
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		RespondErrorJSON(w, "Failed to get gmail message", http.StatusServiceUnavailable, err)
		return
	}

	var rawEmails []models.RawEmail

	for _, message := range messages {
		currentRawEmail := models.RawEmail{
			Sender:    getHeader(message.Payload.Headers, "From"),
			Recipient: getHeader(message.Payload.Headers, "To"),
			Subject:   getHeader(message.Payload.Headers, "Subject"),
		}

		// Skip emails with invalid date
		if message.InternalDate == "" {
			continue
		}
		dateNumber, convErr := strconv.ParseInt(message.InternalDate, 10, 64)
		if convErr != nil {
			RespondErrorJSON(w, "Failed to parse email date", http.StatusServiceUnavailable, convErr)
			return
		}
		date := time.UnixMilli(dateNumber)
		currentRawEmail.Date = date

		decoded, decodeErr := base64.URLEncoding.DecodeString(message.Body.Data)
		if decodeErr != nil {
			RespondErrorJSON(w, "Failed to parse email body", http.StatusServiceUnavailable, decodeErr)
			return
		}
		body := string(decoded)
		currentRawEmail.Body = body

		rawEmails = append(rawEmails, currentRawEmail)
	}

	/**
	 * Trigger email analyzer pipeline
	 */
	message, statusCode, data, err := AnalyzeEmails(h, *user, rawEmails)

	if err != nil {
		RespondErrorJSON(w, message, statusCode, err)
		return
	}

	RespondDataJSON(w, message, statusCode, data)

}

/**
 * Upstream services - gmail
 */

func (h *Handler) getNewAccessToken(refreshToken string) (*TokenRequestResponse, error) {
	/**
	 * Prepare payload (as form data)
	 */
	data := url.Values{}
	data.Set("client_id", h.Config.OauthClientId)
	data.Set("client_secret", h.Config.OauthClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest(
		"POST",
		constants.GoogleOAuthURL,
		strings.NewReader(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	/**
	 * Make request
	 */
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res TokenRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (h *Handler) getGmails(accessToken string) (*GmailsResponse, error) {
	/**
	 * Prepare url with query params
	 */

	baseURL := constants.GmailAPIURL + "/users/me/messages"

	params := url.Values{}
	params.Set("maxResults", constants.MaxGmailCountPerSync)

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	/**
	 * Make request
	 */
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res GmailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (h *Handler) getGmailMessage(id string, accessToken string, ctx context.Context) (*GmailResponse, error) {
	/**
	 * Prepare url with query params
	 */

	baseURL := constants.GmailAPIURL + "/users/me/messages/" + id

	params := url.Values{}
	params.Set("format", "full")

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	/**
	 * Make request
	 */

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res GmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
