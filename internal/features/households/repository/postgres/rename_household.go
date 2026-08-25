package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) RenameHousehold(
	ctx context.Context,
	householdID uuid.UUID,
	name string,
	version int,
) (core_domain.Household, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE households
		SET name = $1, updated_at = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, version, name, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query, name, time.Now(), householdID, version)

	var householdModel HouseholdModel
	err := row.Scan(
		&householdModel.ID,
		&householdModel.Version,
		&householdModel.Name,
		&householdModel.CreatedAt,
		&householdModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.Household{}, fmt.Errorf(
				"household with id %q concurrently accessed: %w",
				householdID, core_errors.ErrConflict,
			)
		}

		return core_domain.Household{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.NewHousehold(
		householdModel.ID,
		householdModel.Version,
		householdModel.Name,
		householdModel.CreatedAt,
		householdModel.UpdatedAt,
	), nil
}
