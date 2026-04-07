package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/thetsGit/spend-wise-be/internal/ai"
	"github.com/thetsGit/spend-wise-be/internal/config"
)

func main() {
	// Bootstrap things and prepare environment
	godotenv.Load()

	config := config.Load()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	r.Get("/sing", func(w http.ResponseWriter, r *http.Request) {
		result, err := ai.CallOpenAI("Sing me a song", config)
		if err != nil {
			fmt.Fprintf(w, `{"status": "error", "error": "%s"}`, err.Error())
			return
		}
		fmt.Fprintf(w, `{"status": "success", "result": "%s"}`, result)
	})

	// fmt.Printf("Server starting on :%s", config)
	http.ListenAndServe(":"+config.HTTPPort, r)
}
