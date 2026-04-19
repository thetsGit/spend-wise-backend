package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/thetsGit/spend-wise-be/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(userContextKey).(*models.User)
	return user
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/**
		 * Extract token from headers
		 */

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			RespondErrorJSON(w, "Unauthorized request", http.StatusUnauthorized, nil)
			return
		}

		sessionToken := strings.Split(authHeader, " ")[1] // (e.g "Bearer xxxx")

		/**
		 * Check DB if the token exists
		 */
		user, err := h.DB.GetUserBySessionToken(sessionToken)

		if err != nil {
			RespondErrorJSON(w, "Invalid session token", http.StatusUnauthorized, nil)
			return
		}

		/**
		 * Check token's expiry
		 */
		if user.ExpiresAt.Before(time.Now()) {
			RespondErrorJSON(w, "Invalid session token", http.StatusUnauthorized, nil)
			return
		}

		// Save to context for later use by other handlers
		ctx := context.WithValue(r.Context(), userContextKey, &user)

		// Pass request to other handlers
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
