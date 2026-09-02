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

func (r *HouseholdsRepository) AcceptInvite(
	ctx context.Context,
	code string,
	userID uuid.UUID,
) (core_domain.Household, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		WITH valid_invite AS (
			SELECT id, household_id
			FROM household_invites
			WHERE code = $1
			  AND revoked_at IS NULL
			  AND expires_at > now()
			  AND (max_uses IS NULL OR use_count < max_uses)
			  AND NOT EXISTS (
			      SELECT 1 FROM household_members
			      WHERE household_id = household_invites.household_id AND user_id = $2
			  )
		),
		bump_invite AS (
			UPDATE household_invites
			SET use_count = use_count + 1
			WHERE id IN (SELECT id FROM valid_invite)
			RETURNING household_id
		),
		new_member AS (
			INSERT INTO household_members (household_id, user_id, role)
			SELECT household_id, $2, 'member' FROM bump_invite
			RETURNING household_id
		)
		SELECT h.id, h.version, h.name, h.created_at, h.updated_at
		FROM households h
		JOIN new_member nm ON nm.household_id = h.id;
	`

	row := r.pool.QueryRow(ctx, query, code, userID)

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
				"invite is invalid, expired, exhausted, or already accepted: %w",
				core_errors.ErrInvalidArgument,
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
