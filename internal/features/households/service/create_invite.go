package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const inviteTTL = 7 * 24 * time.Hour

type CreateInviteRequest struct {
	MaxUses *int `json:"max_uses" validate:"omitempty,min=1"`
}

func generateInviteCode() (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)

	return strings.ToUpper(code), nil
}

func (s *HouseholdsService) CreateInvite(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	request CreateInviteRequest,
) (core_domain.HouseholdInvite, error) {
	caller, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return core_domain.HouseholdInvite{}, fmt.Errorf("get household: %w", err)
	}

	if caller.Role != core_domain.HouseholdRoleOwner && caller.Role != core_domain.HouseholdRoleAdmin {
		return core_domain.HouseholdInvite{}, fmt.Errorf(
			"role %q can't create invites: %w", caller.Role, core_errors.ErrForbidden,
		)
	}

	code, err := generateInviteCode()
	if err != nil {
		return core_domain.HouseholdInvite{}, fmt.Errorf("generate invite code: %w", err)
	}

	invite := core_domain.NewHouseholdInvite(
		householdID,
		code,
		callerID,
		time.Now().Add(inviteTTL),
		request.MaxUses,
	)

	created, err := s.householdsRepository.CreateInvite(ctx, invite)
	if err != nil {
		return core_domain.HouseholdInvite{}, fmt.Errorf("create invite: %w", err)
	}

	return created, nil
}
