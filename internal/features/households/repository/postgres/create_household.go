package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) CreateHousehold(
	ctx context.Context,
	household core_domain.Household,
	ownerID uuid.UUID,
) (core_domain.Household, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH new_household AS (
			INSERT INTO households (version, name, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id, version, name, created_at, updated_at
		),
		new_member AS (
			INSERT INTO household_members (household_id, user_id, role)
			SELECT id, $5, 'owner' FROM new_household
		)
		SELECT id, version, name, created_at, updated_at FROM new_household;
	`

	row := r.pool.QueryRow(ctx, query,
		household.Version,
		household.Name,
		household.CreatedAt,
		household.UpdatedAt,
		ownerID,
	)

	var householdModel HouseholdModel
	err := row.Scan(
		&householdModel.ID,
		&householdModel.Version,
		&householdModel.Name,
		&householdModel.CreatedAt,
		&householdModel.UpdatedAt,
	)
	if err != nil {
		return core_domain.Household{}, fmt.Errorf("scan error: %w", err)
	}

	householdDomain := core_domain.NewHousehold(
		householdModel.ID,
		householdModel.Version,
		householdModel.Name,
		householdModel.CreatedAt,
		householdModel.UpdatedAt,
	)

	return householdDomain, nil
}
