package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) ListLists(
	ctx context.Context,
	householdID uuid.UUID,
) ([]core_domain.ShoppingList, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, household_id, version, name, created_by, created_at, updated_at
		FROM shopping_lists
		WHERE household_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, householdID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result := make([]core_domain.ShoppingList, 0)

	for rows.Next() {
		list, err := scanShoppingList(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, list)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
