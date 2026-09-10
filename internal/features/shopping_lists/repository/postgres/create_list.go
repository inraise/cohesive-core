package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
)

func (r *ShoppingListsRepository) CreateList(
	ctx context.Context,
	list core_domain.ShoppingList,
) (core_domain.ShoppingList, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO shopping_lists (household_id, name, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, household_id, version, name, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		list.HouseholdID,
		list.Name,
		list.CreatedBy,
		list.CreatedAt,
		list.UpdatedAt,
	)

	return scanShoppingList(row)
}
