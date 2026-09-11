package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) GetItemByID(
	ctx context.Context,
	shoppingListID uuid.UUID,
	itemID uuid.UUID,
) (core_domain.ShoppingListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, shopping_list_id, version, name, quantity, is_purchased, created_by, created_at, updated_at
		FROM shopping_list_items
		WHERE id = $1 AND shopping_list_id = $2;
	`

	row := r.pool.QueryRow(ctx, query, itemID, shoppingListID)

	item, err := scanShoppingListItem(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.ShoppingListItem{}, fmt.Errorf(
				"item with id %q in list %q: %w", itemID, shoppingListID, core_errors.ErrNotFound,
			)
		}

		return core_domain.ShoppingListItem{}, err
	}

	return item, nil
}
