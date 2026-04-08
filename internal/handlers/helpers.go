package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/models"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondErrorJSON(w http.ResponseWriter, message string, err error) {
	// Server side logging for obervability
	fmt.Printf("[ERROR] %s: %v\n", message, err)

	respondJSON(w, http.StatusOK, models.APIResponse{
		Status:  "error",
		Message: message,
		Error:   err.Error(),
	})
}

func (h *Handler) respondDataJSON(w http.ResponseWriter, message string, data any) {
	respondJSON(w, http.StatusOK, models.APIResponse{
		Status:  "success",
		Message: "Emails processed",
		Data:    data,
	})
}
