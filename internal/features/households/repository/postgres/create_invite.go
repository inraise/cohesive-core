package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"
)

func (r *HouseholdsRepository) CreateInvite(
	ctx context.Context,
	invite core_domain.HouseholdInvite,
) (core_domain.HouseholdInvite, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO household_invites (household_id, code, created_by, expires_at, max_uses, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, household_id, code, created_by, expires_at, max_uses, use_count, revoked_at, created_at;
	`

	row := r.pool.QueryRow(ctx, query,
		invite.HouseholdID,
		invite.Code,
		invite.CreatedBy,
		invite.ExpiresAt,
		invite.MaxUses,
		invite.CreatedAt,
	)

	var result core_domain.HouseholdInvite
	err := row.Scan(
		&result.ID,
		&result.HouseholdID,
		&result.Code,
		&result.CreatedBy,
		&result.ExpiresAt,
		&result.MaxUses,
		&result.UseCount,
		&result.RevokedAt,
		&result.CreatedAt,
	)
	if err != nil {
		return core_domain.HouseholdInvite{}, fmt.Errorf("scan error: %w", err)
	}

	return result, nil
}
