package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/thetsGit/spend-wise-be/internal/models"
	"github.com/thetsGit/spend-wise-be/internal/utils"
)

type OauthUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type VerifyOauthRequest struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type VerifyOauthResponse struct {
	User         models.PublicUser `json:"user"`
	SessionToken string            `json:"session_token"`
}

func (h *Handler) VerifyOauth(w http.ResponseWriter, r *http.Request) {

	/**
	 * Decode post payload
	 */
	var req VerifyOauthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondErrorJSON(w, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	/**
	 * Fetch user info from google
	 */
	userInfo, err := h.fetchGoogleUserInfo(req.AccessToken)
	if err != nil {
		RespondErrorJSON(w, "Failed to fetch user info from Google", http.StatusUnauthorized, err)
		return
	}

	/**
	 * Issue an Opaque token
	 */
	sessionToken := utils.GenerateOpaqueToken()
	expiresAt := time.Now().UTC().Add(h.Config.AuthSessionLifeSec)

	userToSave := models.User{
		SessionToken: sessionToken,
		ExpiresAt:    expiresAt,
		PublicUser: models.PublicUser{
			OauthId:      userInfo.ID,
			OauthEmail:   userInfo.Email,
			OauthName:    userInfo.Name,
			OauthPicture: &userInfo.Picture,
		},

		OauthAccessToken:  req.AccessToken,
		OauthRefreshToken: req.RefreshToken,
		OauthTokenExpiry:  time.Now().UTC().Add(time.Duration(req.ExpiresIn) * time.Second),
		OauthTokenType:    req.TokenType,
		OauthScope:        req.Scope,
	}

	fmt.Println(sessionToken, expiresAt, userInfo, userToSave)

	/**
	 * Save or update user info to 'users' DB
	 */
	savedUser, err := h.DB.UpsertUser(userToSave)

	if err != nil {
		RespondErrorJSON(w, "Failed to save user information", http.StatusUnauthorized, err)
		return
	}

	RespondDataJSON(w, "Success", http.StatusOK, VerifyOauthResponse{
		User:         savedUser.PublicUser,
		SessionToken: savedUser.SessionToken,
	})
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Potentially fetched user during auth middleware check
	user := GetUserFromContext(r.Context())

	if user == nil {
		RespondErrorJSON(w, "Failed to get user information", http.StatusBadRequest, nil)
		return
	}

	RespondDataJSON(w, "Success", http.StatusOK, user.PublicUser)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())

	if user == nil {
		RespondErrorJSON(w, "User not found", http.StatusBadRequest, nil)
		return
	}

	/**
	 * Revoke token
	 */

	err := h.DB.ClearUserSession(user.SessionToken)

	if err != nil {
		RespondErrorJSON(w, "Failed to log out", http.StatusBadRequest, err)
		return
	}

	RespondDataJSON(w, "Success", http.StatusNoContent, nil)
}

/**
 * Upstream services
 */

func (h *Handler) fetchGoogleUserInfo(accessToken string) (*OauthUserInfo, error) {
	req, _ := http.NewRequest("GET", h.Config.OauthApiUrl+"/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo OauthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
