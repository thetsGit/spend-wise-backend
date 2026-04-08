package handlers

import (
	"net/http"

	"github.com/thetsGit/spend-wise-be/internal/models"
)

func (h *Handler) GetSaasDiscoveries(w http.ResponseWriter, r *http.Request) {
	filter := models.SaaSDiscoveryFilter{
		ProductName: r.URL.Query().Get("product_name"),
		SignalType:  r.URL.Query().Get("signal_type"),
	}

	results, err := h.DB.GetSaaSDiscoveries(filter)
	if err != nil {
		RespondErrorJSON(w, "Failed to fetch saas discovery list", http.StatusInternalServerError, err)
		return
	}

	RespondDataJSON(w, "Success", http.StatusOK, results)
}

func (h *Handler) GetSaasDiscoverySummary(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.GetSaaSDiscoverySummary()
	if err != nil {
		RespondErrorJSON(w, "Failed to fetch saas discovery summary", http.StatusInternalServerError, err)
		return
	}

	RespondDataJSON(w, "Success", http.StatusOK, results)
}
