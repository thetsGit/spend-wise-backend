package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/models"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondErrorJSON(w http.ResponseWriter, message string, err error) {
	// Server side logging for obervability
	fmt.Printf("[ERROR] %s: %v\n", message, err)

	parsedError := ""
	if err != nil {
		parsedError = err.Error()
	}
	RespondJSON(w, http.StatusOK, models.APIResponse{
		Status:  "error",
		Message: message,
		Error:   parsedError,
	})

}

func RespondDataJSON(w http.ResponseWriter, message string, data any) {
	RespondJSON(w, http.StatusOK, models.APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}
