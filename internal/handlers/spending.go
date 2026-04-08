package handlers

import (
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/models"
)

func (h *Handler) GetSpending(w http.ResponseWriter, r *http.Request) {
	filter := models.SpendingFilter{
		Category:  r.URL.Query().Get("category"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	results, err := h.DB.GetSpending(filter)
	if err != nil {
		RespondErrorJSON(w, "Failed to fetch spending list", http.StatusInternalServerError, err)
		return
	}

	RespondDataJSON(w, "Success", http.StatusOK, results)
}

func (h *Handler) GetSpendingSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.DB.GetSpendingSummary()
	if err != nil {
		RespondErrorJSON(w, "Failed to fetch spending summary", http.StatusInternalServerError, err)
		return
	}

	RespondDataJSON(w, "Success", http.StatusOK, summary)
}
