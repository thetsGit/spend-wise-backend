package handlers

import (
	"github.com/thetsGit/spend-wise-be/internal/config"
	"github.com/thetsGit/spend-wise-be/internal/database"
)

type Handler struct {
	DB     *database.DB
	Config *config.Config
}

func CreateHandlers(db *database.DB, config *config.Config) Handler {
	return Handler{
		DB:     db,
		Config: config,
	}
}
