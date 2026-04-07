package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/thetsGit/spend-wise-be/internal/config"
	"github.com/thetsGit/spend-wise-be/internal/database"
)

func hello(w http.ResponseWriter, req *http.Request) {
	config := config.Load()
	db, err := database.Connect(config)

	if err != nil {
		fmt.Fprintf(w, "failed to create pool: %w", err)
		return
	}

	db.Close()

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	// Bootstrap things and prepare environment
	godotenv.Load()

	config := config.Load()

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(fmt.Sprintf(":%v", config.HTTP_PORT), nil)
}
