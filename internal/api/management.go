package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"sledgehammer.echo-mesh.com/internal/cryptography"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	databasePasswordHash, role, err := a.Store.GetPasswordHashAndRole(credentials.Username)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	hashedPassword, err := cryptography.HashPassword(credentials.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	if !cryptography.CheckPasswordHash(hashedPassword, databasePasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expiration := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Username: credentials.Username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "could not login", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (a *API) Logout() {}

func (a *API) ApproveBan() {}

func (a *API) RejectBan() {}

func (a *API) DeleteBan() {}

func (a *API) GetBanHistory() {}

func (a *API) GetWaitingForApprovalBans() {}
