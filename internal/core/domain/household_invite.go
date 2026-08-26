package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type HouseholdInvite struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Code        string
	CreatedBy   uuid.UUID

	ExpiresAt time.Time
	MaxUses   *int
	UseCount  int
	RevokedAt *time.Time
	CreatedAt time.Time
}

func NewHouseholdInvite(
	householdID uuid.UUID,
	code string,
	createdBy uuid.UUID,
	expiresAt time.Time,
	maxUses *int,
) HouseholdInvite {
	return HouseholdInvite{
		ID:          UninitializedID,
		HouseholdID: householdID,
		Code:        code,
		CreatedBy:   createdBy,
		ExpiresAt:   expiresAt,
		MaxUses:     maxUses,
		UseCount:    0,
		CreatedAt:   time.Now(),
	}
}
