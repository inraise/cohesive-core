package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) GetHousehold(
	ctx context.Context,
	householdID uuid.UUID,
	userID uuid.UUID,
) (core_domain.HouseholdWithRole, error) {
	household, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, userID)
	if err != nil {
		return core_domain.HouseholdWithRole{}, fmt.Errorf("get household from repository: %w", err)
	}

	return household, nil
}
