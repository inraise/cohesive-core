package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) ListInvites(
	ctx context.Context,
	householdID uuid.UUID,
) ([]core_domain.HouseholdInvite, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, household_id, code, created_by, expires_at, max_uses, use_count, revoked_at, created_at
		FROM household_invites
		WHERE household_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, householdID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result := make([]core_domain.HouseholdInvite, 0)

	for rows.Next() {
		var invite core_domain.HouseholdInvite

		if err := rows.Scan(
			&invite.ID,
			&invite.HouseholdID,
			&invite.Code,
			&invite.CreatedBy,
			&invite.ExpiresAt,
			&invite.MaxUses,
			&invite.UseCount,
			&invite.RevokedAt,
			&invite.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		result = append(result, invite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
