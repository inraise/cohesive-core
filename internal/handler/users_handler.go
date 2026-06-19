package handler

import (
	"encoding/json"
	"net/http"

	"cohesive-core/internal/models"
	"cohesive-core/internal/service"
)

type AuthHandler struct {
	services *service.AuthService
}

func NewAuthHandler(services *service.AuthService) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) RegistrationUser(w http.ResponseWriter, r *http.Request) {
	var input models.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и Password обязательны", http.StatusBadRequest)
		return
	}

	token, err := h.services.RegistrationUser(r.Context(), input)
	if err != nil {
		http.Error(w, "Ошибка регистрации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var input models.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email и Password обязательны", http.StatusBadRequest)
		return
	}

	token, err := h.services.LoginUser(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}