package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) DeleteHousehold(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) error {
	current, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return fmt.Errorf("get household: %w", err)
	}

	if current.Role != core_domain.HouseholdRoleOwner {
		return fmt.Errorf("only owner can delete household: %w", core_errors.ErrForbidden)
	}

	if err := s.householdsRepository.DeleteHousehold(ctx, householdID); err != nil {
		return fmt.Errorf("delete household: %w", err)
	}

	return nil
}
