package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type HouseholdRole string

const (
	HouseholdRoleOwner  HouseholdRole = "owner"
	HouseholdRoleAdmin  HouseholdRole = "admin"
	HouseholdRoleMember HouseholdRole = "member"
)

type HouseholdWithRole struct {
	Household Household
	Role      HouseholdRole
}

type HouseholdMember struct {
	UserID    uuid.UUID
	Email     string
	FirstName string
	LastName  *string

	Role     HouseholdRole
	JoinedAt time.Time
}
