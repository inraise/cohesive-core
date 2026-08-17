package households_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type HouseholdModel struct {
	ID      uuid.UUID
	Version int

	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}
