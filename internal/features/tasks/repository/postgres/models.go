package tasks_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type TaskModel struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Version     int

	Title       string
	Description *string
	Status      string
	AssignedTo  *uuid.UUID
	CreatedBy   uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}
