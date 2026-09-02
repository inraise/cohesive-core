package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) RevokeInvite(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	inviteID uuid.UUID,
) error {
	caller, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return fmt.Errorf("get household: %w", err)
	}

	if caller.Role != core_domain.HouseholdRoleOwner && caller.Role != core_domain.HouseholdRoleAdmin {
		return fmt.Errorf("role %q can't revoke invites: %w", caller.Role, core_errors.ErrForbidden)
	}

	if err := s.householdsRepository.RevokeInvite(ctx, householdID, inviteID); err != nil {
		return fmt.Errorf("revoke invite: %w", err)
	}

	return nil
}
