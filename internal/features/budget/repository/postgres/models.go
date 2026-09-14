package budget_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Version     int

	Type        string
	Amount      int64
	Description *string
	CreatedBy   uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}
