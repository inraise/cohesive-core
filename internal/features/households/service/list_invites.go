package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) ListInvites(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) ([]core_domain.HouseholdInvite, error) {
	caller, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return nil, fmt.Errorf("get household: %w", err)
	}

	if caller.Role != core_domain.HouseholdRoleOwner && caller.Role != core_domain.HouseholdRoleAdmin {
		return nil, fmt.Errorf("role %q can't view invites: %w", caller.Role, core_errors.ErrForbidden)
	}

	invites, err := s.householdsRepository.ListInvites(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("list invites from repository: %w", err)
	}

	return invites, nil
}
