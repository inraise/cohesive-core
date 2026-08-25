package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RenameHouseholdRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (s *HouseholdsService) RenameHousehold(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	request RenameHouseholdRequest,
) (core_domain.HouseholdWithRole, error) {
	current, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID)
	if err != nil {
		return core_domain.HouseholdWithRole{}, fmt.Errorf("get household: %w", err)
	}

	if current.Role != core_domain.HouseholdRoleOwner && current.Role != core_domain.HouseholdRoleAdmin {
		return core_domain.HouseholdWithRole{}, fmt.Errorf(
			"role %q can't rename household: %w", current.Role, core_errors.ErrForbidden,
		)
	}

	updated, err := s.householdsRepository.RenameHousehold(ctx, householdID, request.Name, current.Household.Version)
	if err != nil {
		return core_domain.HouseholdWithRole{}, fmt.Errorf("rename household: %w", err)
	}

	return core_domain.HouseholdWithRole{Household: updated, Role: current.Role}, nil
}
