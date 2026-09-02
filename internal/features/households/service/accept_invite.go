package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HouseholdsService) AcceptInvite(
	ctx context.Context,
	code string,
	userID uuid.UUID,
) (core_domain.Household, error) {
	household, err := s.householdsRepository.AcceptInvite(ctx, code, userID)
	if err != nil {
		return core_domain.Household{}, fmt.Errorf("accept invite: %w", err)
	}

	return household, nil
}
