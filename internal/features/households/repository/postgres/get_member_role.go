package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) GetMemberRole(
	ctx context.Context,
	householdID uuid.UUID,
	userID uuid.UUID,
) (core_domain.HouseholdRole, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT role FROM household_members
		WHERE household_id = $1 AND user_id = $2;
	`

	row := r.pool.QueryRow(ctx, query, householdID, userID)

	var role string
	if err := row.Scan(&role); err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return "", fmt.Errorf(
				"user %q is not a member of household %q: %w",
				userID, householdID, core_errors.ErrNotFound,
			)
		}

		return "", fmt.Errorf("scan error: %w", err)
	}

	return core_domain.HouseholdRole(role), nil
}
