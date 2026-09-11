package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
)

func (r *ShoppingListsRepository) CreateItem(
	ctx context.Context,
	item core_domain.ShoppingListItem,
) (core_domain.ShoppingListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO shopping_list_items (shopping_list_id, name, quantity, is_purchased, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, shopping_list_id, version, name, quantity, is_purchased, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		item.ShoppingListID,
		item.Name,
		item.Quantity,
		item.IsPurchased,
		item.CreatedBy,
		item.CreatedAt,
		item.UpdatedAt,
	)

	return scanShoppingListItem(row)
}
