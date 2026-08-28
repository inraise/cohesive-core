package households_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) TransferOwnership(
	ctx context.Context,
	householdID uuid.UUID,
	currentOwnerID uuid.UUID,
	newOwnerID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH demoted AS (
			UPDATE household_members
			SET role = 'admin'
			WHERE household_id = $1 AND user_id = $2 AND role = 'owner'
			RETURNING household_id
		)
		UPDATE household_members
		SET role = 'owner'
		FROM demoted
		WHERE household_members.household_id = demoted.household_id
		  AND household_members.user_id = $3
		RETURNING household_members.user_id;
	`

	row := r.pool.QueryRow(ctx, query, householdID, currentOwnerID, newOwnerID)

	var returnedID uuid.UUID
	if err := row.Scan(&returnedID); err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return fmt.Errorf(
				"caller is not the current owner or target is not a member: %w",
				core_errors.ErrConflict,
			)
		}

		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}
