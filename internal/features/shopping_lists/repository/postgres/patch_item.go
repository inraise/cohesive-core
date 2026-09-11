package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *ShoppingListsRepository) PatchItem(
	ctx context.Context,
	item core_domain.ShoppingListItem,
) (core_domain.ShoppingListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE shopping_list_items
		SET name = $1, quantity = $2, is_purchased = $3, updated_at = $4, version = version + 1
		WHERE id = $5 AND shopping_list_id = $6 AND version = $7
		RETURNING id, shopping_list_id, version, name, quantity, is_purchased, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		item.Name,
		item.Quantity,
		item.IsPurchased,
		time.Now(),
		item.ID,
		item.ShoppingListID,
		item.Version,
	)

	updated, err := scanShoppingListItem(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.ShoppingListItem{}, fmt.Errorf(
				"item with id %q concurrently accessed: %w", item.ID, core_errors.ErrConflict,
			)
		}

		return core_domain.ShoppingListItem{}, err
	}

	return updated, nil
}
