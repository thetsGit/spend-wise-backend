package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/models"
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

	var rawEmails []models.RawEmail
	err := json.NewDecoder(r.Body).Decode(&rawEmails)
	if err != nil {
		RespondErrorJSON(w, "Request body too large or invalid JSON", http.StatusBadRequest, err)
		return
	}

	user := GetUserFromContext(r.Context())
	message, statusCode, data, err := AnalyzeEmails(h, *user, rawEmails)

	if err != nil {
		RespondErrorJSON(w, message, statusCode, err)
		return
	}

	RespondDataJSON(w, message, statusCode, data)
}
