package handler

import (
	"cohesive-core/internal/models"
	"cohesive-core/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type FamilyHandler struct {
	service *service.FamilyService
}

func NewFamilyHandler(service *service.FamilyService) *FamilyHandler {
	return &FamilyHandler{
		service: service,
	}
}

func (h *FamilyHandler) CreateFamily(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("UserId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Отсутствует или неверный заголовок UserId", http.StatusUnauthorized)
		return
	}

	var input models.CreateFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "Имя семьи обязательно", http.StatusBadRequest)
		return
	}

	family, err := h.service.CreateFamily(r.Context(), input, userID)
	if err != nil {
		http.Error(w, "Ошибка создания семьи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(family)
}
