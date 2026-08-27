package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) RemoveMember(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	targetID uuid.UUID,
) error {
	caller, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return fmt.Errorf("get household: %w", err)
	}

	if targetID == callerID {
		if caller.Role == core_domain.HouseholdRoleOwner {
			return fmt.Errorf(
				"owner must transfer ownership or delete household before leaving: %w",
				core_errors.ErrConflict,
			)
		}

		if err := s.householdsRepository.RemoveMember(ctx, householdID, targetID); err != nil {
			return fmt.Errorf("remove member: %w", err)
		}

		return nil
	}

	targetRole, err := s.householdsRepository.GetMemberRole(ctx, householdID, targetID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return err
		}

		return fmt.Errorf("get target member role: %w", err)
	}

	switch caller.Role {
	case core_domain.HouseholdRoleOwner:
	case core_domain.HouseholdRoleAdmin:
		if targetRole != core_domain.HouseholdRoleMember {
			return fmt.Errorf(
				"admin can only remove regular members, not other admins or the owner: %w",
				core_errors.ErrForbidden,
			)
		}
	default:
		return fmt.Errorf("role %q can't remove members: %w", caller.Role, core_errors.ErrForbidden)
	}

	if err := s.householdsRepository.RemoveMember(ctx, householdID, targetID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}

	return nil
}
