package handler

import (
	"cohesive-core/internal/models"
	"cohesive-core/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
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

	var input models.FamilyRequest
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

func (h *FamilyHandler) UpdateFamily(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("UserId")
	familyIDStr := r.Header.Get("FamilyId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Отсутствует или неверный заголовок UserId", http.StatusUnauthorized)
		return
	}

	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		http.Error(w, "Отсутствует или неверный заголовок FamilyId", http.StatusBadRequest)
		return
	}

	var input models.FamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateFamily(r.Context(), familyID, userID, input)
	if err != nil {
		if err.Error() == "У вас нет прав на редактирование этой семьи или она не существует" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func (h *FamilyHandler) JoinFamily(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserId")

	var input models.JoinFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.JoinFamily(r.Context(), input, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "joined"}`))
}

func (h *FamilyHandler) LeaveFamily(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserId")

	if err := h.service.LeaveFamily(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "left"}`))
}

func (h *FamilyHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	actorID := r.Header.Get("UserId")
	familyID := r.Header.Get("FamilyId")

	fmt.Printf("Дебаг PATCH: ActorID=%s, FamilyID=%s\n", actorID, familyID)

	targetUserStr := chi.URLParam(r, "user_id")
	_, err := uuid.Parse(targetUserStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя в URL", http.StatusBadRequest)
		return
	}

	var input models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateMemberRole(r.Context(), familyID, actorID, targetUserStr, input.Role); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "role updated"}`))
}

func (h *FamilyHandler) KickMember(w http.ResponseWriter, r *http.Request) {
	actorID := r.Header.Get("UserId")
	familyID := r.Header.Get("FamilyId")

	targetUserStr := chi.URLParam(r, "user_id")
	_, err := uuid.Parse(targetUserStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя в URL", http.StatusBadRequest)
		return
	}

	if err := h.service.KickMember(r.Context(), familyID, actorID, targetUserStr); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "member kicked"}`))
}
