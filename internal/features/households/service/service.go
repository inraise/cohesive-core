package households_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"

	"github.com/google/uuid"
)

type HouseholdsService struct {
	householdsRepository HouseholdsRepository
}

type HouseholdsRepository interface {
	CreateHousehold(
		ctx context.Context,
		household core_domain.Household,
		ownerID uuid.UUID,
	) (core_domain.Household, error)
}

func NewHouseholdsService(
	householdsRepository HouseholdsRepository,
) *HouseholdsService {
	return &HouseholdsService{
		householdsRepository: householdsRepository,
	}
}
