package households_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) RevokeInvite(
	ctx context.Context,
	householdID uuid.UUID,
	inviteID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE household_invites
		SET revoked_at = $1
		WHERE id = $2 AND household_id = $3 AND revoked_at IS NULL;
	`

	tag, err := r.pool.Exec(ctx, query, time.Now(), inviteID, householdID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"invite not found or already revoked: %w",
			core_errors.ErrNotFound,
		)
	}

	return nil
}
