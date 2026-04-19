package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	/**
	 * API services
	 */

	r.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   config.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
			AllowCredentials: true,
		}))

		/**
		 * Public routes
		 */

		r.Post("/oauth/verify", handler.VerifyOauth)

		/**
		 * Protected routes
		 */

		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)

			r.Post("/emails/upload", handler.UploadEmails)

			r.Get("/spending", handler.GetSpending)
			r.Get("/spending/summary", handler.GetSpendingSummary)

			r.Get("/saas", handler.GetSaasDiscoveries)
			r.Get("/saas/summary", handler.GetSaasDiscoverySummary)

			r.Get("/users/me", handler.GetMe)

			r.Post("/auth/logout", handler.Logout)
		})

	})

	/**
	 * Fallback / error routes
	 */

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondErrorJSON(w, "Route not found", http.StatusNotFound, nil)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
	})

	/**
	 * Start the server
	 */

	fmt.Printf("Server starting on :%s", config)
	http.ListenAndServe(":"+config.HTTPPort, r)
}
