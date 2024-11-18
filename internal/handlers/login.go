package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/models"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/services"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/validators"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler godoc
// @Summary Log in a user
// @Description Log in a user by validating email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body models.LoginRequest true "Login Request"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} map[string]string "Validation errors"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validators.GetValidator().Struct(creds); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(creds.Email)
	if err != nil || user.PasswordHash != nil || bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := services.GenerateToken(user.Email, secret)
	json.NewEncoder(w).Encode(models.TokenResponse{Token: token})
}
