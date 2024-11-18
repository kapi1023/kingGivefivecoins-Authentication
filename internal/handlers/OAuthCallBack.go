package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/oauth"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/services"
)

// OAuthCallbackHandler godoc
// @Summary Handle OAuth callback
// @Description Handle OAuth callback for Google, LinkedIn, and Apple
// @Tags Authentication
// @Accept json
// @Produce json
// @Param provider query string true "OAuth provider (google, linkedin, apple)"
// @Param code query string true "Authorization code"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /oauth/callback [get]
func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	code := r.URL.Query().Get("code")

	if provider == "" || code == "" {
		http.Error(w, "Missing provider or code", http.StatusBadRequest)
		return
	}

	// Map provider to respective credentials
	clientID, clientSecret, redirectURL := "", "", "http://localhost:8080/oauth/callback?provider="+provider
	switch provider {
	case string(oauth.Google):
		clientID = cfg.GoogleClient
		clientSecret = cfg.GoogleSecret
	case string(oauth.LinkedIn):
		clientID = cfg.LinkedInClient
		clientSecret = cfg.LinkedInSecret
	case string(oauth.Apple):
		clientID = cfg.AppleClient
		clientSecret = cfg.AppleSecret
	default:
		http.Error(w, "Unsupported provider", http.StatusBadRequest)
		return
	}

	token, err := oauth.ExchangeCodeForToken(context.Background(), oauth.OAuthProvider(provider), code, clientID, clientSecret, redirectURL)
	if err != nil {
		http.Error(w, "Error exchanging code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, err := oauth.ExtractEmailFromToken(token, oauth.OAuthProvider(provider))
	if err != nil {
		http.Error(w, "Error extracting email from token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := db.GetOrCreateUserByOAuth(email, provider)
	if err != nil {
		http.Error(w, "Error creating or retrieving user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Error creating or retrieving user", http.StatusInternalServerError)
		return
	}

	tokenStr, err := services.GenerateToken(user.Email, string(cfg.JWTKey))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token
	response := map[string]string{"token": tokenStr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
