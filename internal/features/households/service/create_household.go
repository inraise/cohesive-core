package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateHouseholdRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (s *HouseholdsService) CreateHousehold(
	ctx context.Context,
	ownerID uuid.UUID,
	request CreateHouseholdRequest,
) (core_domain.Household, error) {
	householdDomain := core_domain.NewHouseholdUninitialized(request.Name)

	if err := householdDomain.Validate(); err != nil {
		return core_domain.Household{}, fmt.Errorf("validate household domain: %w", err)
	}

	household, err := s.householdsRepository.CreateHousehold(ctx, householdDomain, ownerID)
	if err != nil {
		return core_domain.Household{}, fmt.Errorf("create household: %w", err)
	}

	return household, nil
}
