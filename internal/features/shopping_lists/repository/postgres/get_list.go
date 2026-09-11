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

func (r *ShoppingListsRepository) GetListByID(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
) (core_domain.ShoppingList, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, household_id, version, name, created_by, created_at, updated_at
		FROM shopping_lists
		WHERE id = $1 AND household_id = $2;
	`

	row := r.pool.QueryRow(ctx, query, listID, householdID)

	list, err := scanShoppingList(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.ShoppingList{}, fmt.Errorf(
				"shopping list with id %q in household %q: %w", listID, householdID, core_errors.ErrNotFound,
			)
		}

		return core_domain.ShoppingList{}, err
	}

	return list, nil
}
