package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/models"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/services"
)

// RefreshTokenHandler godoc
// @Summary Refresh user token
// @Description Refresh JWT token using the old token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body models.TokenRefreshRequest true "Token Refresh Request"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} map[string]string "Validation errors"
// @Failure 401 {string} string "Invalid token"
// @Failure 500 {string} string "Internal server error"
// @Router /refresh [post]
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request models.TokenRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	claims, err := services.ValidateToken(request.Token, secret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	newToken, _ := services.GenerateToken(claims.Email, secret)
	json.NewEncoder(w).Encode(models.TokenResponse{Token: newToken})
}
