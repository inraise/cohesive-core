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

type MemberDTOResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

func memberDTOFromDomain(member core_domain.HouseholdMember) MemberDTOResponse {
	return MemberDTOResponse{
		UserID:    member.UserID,
		Email:     member.Email,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Role:      string(member.Role),
		JoinedAt:  member.JoinedAt,
	}
}

type InviteDTOResponse struct {
	ID        uuid.UUID  `json:"id"`
	Code      string     `json:"code"`
	ExpiresAt time.Time  `json:"expires_at"`
	MaxUses   *int       `json:"max_uses"`
	UseCount  int        `json:"use_count"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func inviteDTOFromDomain(invite core_domain.HouseholdInvite) InviteDTOResponse {
	return InviteDTOResponse{
		ID:        invite.ID,
		Code:      invite.Code,
		ExpiresAt: invite.ExpiresAt,
		MaxUses:   invite.MaxUses,
		UseCount:  invite.UseCount,
		RevokedAt: invite.RevokedAt,
		CreatedAt: invite.CreatedAt,
	}
}
