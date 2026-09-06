package tasks_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type TaskDTOResponse struct {
	ID          uuid.UUID `json:"id"`
	HouseholdID uuid.UUID `json:"household_id"`
	Version     int       `json:"version"`

	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
	CreatedBy   uuid.UUID  `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func taskDTOFromDomain(task core_domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:          task.ID,
		HouseholdID: task.HouseholdID,
		Version:     task.Version,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		AssignedTo:  task.AssignedTo,
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
