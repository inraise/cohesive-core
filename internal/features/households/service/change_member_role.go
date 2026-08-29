package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ChangeMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=owner admin member"`
}

func (s *HouseholdsService) ChangeMemberRole(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	targetID uuid.UUID,
	request ChangeMemberRoleRequest,
) error {
	if targetID == callerID {
		return fmt.Errorf(
			"can't change your own role through this endpoint: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	caller, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return fmt.Errorf("get household: %w", err)
	}

	if caller.Role != core_domain.HouseholdRoleOwner {
		return fmt.Errorf("only owner can change member roles: %w", core_errors.ErrForbidden)
	}

	newRole := core_domain.HouseholdRole(request.Role)

	if newRole == core_domain.HouseholdRoleOwner {
		if err := s.householdsRepository.TransferOwnership(ctx, householdID, callerID, targetID); err != nil {
			return fmt.Errorf("transfer ownership: %w", err)
		}

		return nil
	}

	if err := s.householdsRepository.SetMemberRole(ctx, householdID, targetID, request.Role); err != nil {
		return fmt.Errorf("set member role: %w", err)
	}

	return nil
}
