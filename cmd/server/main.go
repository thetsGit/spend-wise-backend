package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/thetsGit/spend-wise-be/internal/config"
	"github.com/thetsGit/spend-wise-be/internal/database"
	"github.com/thetsGit/spend-wise-be/internal/handlers"
)

func main() {
	/**
	 * Bootstrap required things (e.g env vars) and make environment get ready
	 */

	godotenv.Load()

	config := config.Load()

	connection, err := database.Connect(config)

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Route handler instance
	handler := handlers.CreateHandlers(connection, config)

	r := chi.NewRouter()

	r.Post("/api/emails/upload", handler.UploadEmails)
	r.Get("/api/spending", handler.GetSpending)
	r.Get("/api/spending/summary", handler.GetSpendingSummary)
	r.Get("/api/saas", handler.GetSaasDiscoveries)
	r.Get("/api/saas/summary", handler.GetSaasDiscoverySummary)

	/**
	 * Fallback / error routes
	 */

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondErrorJSON(w, "Route not found", http.StatusNotFound, nil)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
	})

	fmt.Printf("Server starting on :%s", config)
	http.ListenAndServe(":"+config.HTTPPort, r)
}
