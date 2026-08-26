package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) ListMembers(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) ([]core_domain.HouseholdMember, error) {
	if _, err := s.householdsRepository.GetHouseholdByIDForUser(ctx, householdID, callerID); err != nil {
		return nil, fmt.Errorf("get household: %w", err)
	}

	members, err := s.householdsRepository.ListMembers(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("list members from repository: %w", err)
	}

	return members, nil
}
