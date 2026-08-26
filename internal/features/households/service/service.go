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

	ListHouseholdsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]core_domain.HouseholdWithRole, error)

	GetHouseholdByIDForUser(
		ctx context.Context,
		householdID uuid.UUID,
		userID uuid.UUID,
	) (core_domain.HouseholdWithRole, error)

	RenameHousehold(
		ctx context.Context,
		householdID uuid.UUID,
		name string,
		version int,
	) (core_domain.Household, error)

	DeleteHousehold(
		ctx context.Context,
		householdID uuid.UUID,
	) error

	ListMembers(
		ctx context.Context,
		householdID uuid.UUID,
	) ([]core_domain.HouseholdMember, error)
}

func NewHouseholdsService(
	householdsRepository HouseholdsRepository,
) *HouseholdsService {
	return &HouseholdsService{
		householdsRepository: householdsRepository,
	}
}
