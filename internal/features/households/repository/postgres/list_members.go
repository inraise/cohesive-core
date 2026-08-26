package households_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) ListMembers(
	ctx context.Context,
	householdID uuid.UUID,
) ([]core_domain.HouseholdMember, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, hm.role, hm.joined_at
		FROM household_members hm
		JOIN users u ON u.id = hm.user_id
		WHERE hm.household_id = $1
		ORDER BY hm.joined_at ASC;
	`

	rows, err := r.pool.Query(ctx, query, householdID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result := make([]core_domain.HouseholdMember, 0)

	for rows.Next() {
		var (
			member core_domain.HouseholdMember
			role   string
		)

		if err := rows.Scan(
			&member.UserID,
			&member.Email,
			&member.FirstName,
			&member.LastName,
			&role,
			&member.JoinedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		member.Role = core_domain.HouseholdRole(role)
		result = append(result, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
