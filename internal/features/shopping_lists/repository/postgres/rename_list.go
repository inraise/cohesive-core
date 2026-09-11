package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) RenameList(
	ctx context.Context,
	listID uuid.UUID,
	name string,
	version int,
) (core_domain.ShoppingList, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE shopping_lists
		SET name = $1, updated_at = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, household_id, version, name, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query, name, time.Now(), listID, version)

	list, err := scanShoppingList(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.ShoppingList{}, fmt.Errorf(
				"shopping list with id %q concurrently accessed: %w", listID, core_errors.ErrConflict,
			)
		}

		return core_domain.ShoppingList{}, err
	}

	return list, nil
}
