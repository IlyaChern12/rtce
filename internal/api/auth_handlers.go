package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/IlyaChern12/rtce/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

type HandlerRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// хэндлер для регистрации
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request HandlerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := ah.auth.Register(r.Context(), request.Email, request.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "user already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "user created"})
}

// хэндлер для авторизации
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request HandlerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := ah.auth.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}