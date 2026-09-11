package shoppinglists_repository_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) IsHouseholdMember(
	ctx context.Context,
	householdID uuid.UUID,
	userID uuid.UUID,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM household_members WHERE household_id = $1 AND user_id = $2);`

	row := r.pool.QueryRow(ctx, query, householdID, userID)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("scan error: %w", err)
	}

	return exists, nil
}
