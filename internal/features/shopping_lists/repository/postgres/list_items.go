package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) ListItems(
	ctx context.Context,
	shoppingListID uuid.UUID,
) ([]core_domain.ShoppingListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, shopping_list_id, version, name, quantity, is_purchased, created_by, created_at, updated_at
		FROM shopping_list_items
		WHERE shopping_list_id = $1
		ORDER BY created_at ASC;
	`

	rows, err := r.pool.Query(ctx, query, shoppingListID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result := make([]core_domain.ShoppingListItem, 0)

	for rows.Next() {
		item, err := scanShoppingListItem(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
