package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) ListMyHouseholds(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.HouseholdWithRole, error) {
	households, err := s.householdsRepository.ListHouseholdsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list households from repository: %w", err)
	}

	return households, nil
}
