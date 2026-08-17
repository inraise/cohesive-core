package households_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type HouseholdDTOResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	Name string `json:"name"`
	Role string `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func householdDTOFromDomain(household core_domain.Household, role core_domain.HouseholdRole) HouseholdDTOResponse {
	return HouseholdDTOResponse{
		ID:        household.ID,
		Version:   household.Version,
		Name:      household.Name,
		Role:      string(role),
		CreatedAt: household.CreatedAt,
		UpdatedAt: household.UpdatedAt,
	}
}
