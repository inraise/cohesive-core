package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) ListHouseholdsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.HouseholdWithRole, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT h.id, h.version, h.name, h.created_at, h.updated_at, hm.role
		FROM households h
		JOIN household_members hm ON hm.household_id = h.id
		WHERE hm.user_id = $1
		ORDER BY h.created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result := make([]core_domain.HouseholdWithRole, 0)

	for rows.Next() {
		var (
			householdModel HouseholdModel
			role           string
		)

		if err := rows.Scan(
			&householdModel.ID,
			&householdModel.Version,
			&householdModel.Name,
			&householdModel.CreatedAt,
			&householdModel.UpdatedAt,
			&role,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		result = append(result, core_domain.HouseholdWithRole{
			Household: core_domain.NewHousehold(
				householdModel.ID,
				householdModel.Version,
				householdModel.Name,
				householdModel.CreatedAt,
				householdModel.UpdatedAt,
			),
			Role: core_domain.HouseholdRole(role),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
