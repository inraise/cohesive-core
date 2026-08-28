package households_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) SetMemberRole(
	ctx context.Context,
	householdID uuid.UUID,
	userID uuid.UUID,
	role string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE household_members
		SET role = $1
		WHERE household_id = $2 AND user_id = $3 AND role != 'owner';
	`

	tag, err := r.pool.Exec(ctx, query, role, householdID, userID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"member not found or is the owner: %w",
			core_errors.ErrNotFound,
		)
	}

	return nil
}
