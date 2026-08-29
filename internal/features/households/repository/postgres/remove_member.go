package households_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) RemoveMember(
	ctx context.Context,
	householdID uuid.UUID,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM household_members WHERE household_id = $1 AND user_id = $2;`

	tag, err := r.pool.Exec(ctx, query, householdID, userID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user %q is not a member of household %q: %w",
			userID, householdID, core_errors.ErrNotFound,
		)
	}

	return nil
}
