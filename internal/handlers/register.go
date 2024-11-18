package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/models"
	"github.com/kapi1023/kingGivefivecoins-Authentication/internal/validators"
	"golang.org/x/crypto/bcrypt"
)

var creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user by providing an email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body models.RegisterRequest true "Register Request"
// @Success 201 {string} string "User successfully registered"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validators.GetValidator().Struct(creds); err != nil {
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err := db.CreateUser(creds.Email, string(hashedPassword)); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
