package shoppinglists_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) DeleteList(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM shopping_lists WHERE id = $1 AND household_id = $2;`

	tag, err := r.pool.Exec(ctx, query, listID, householdID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("shopping list with id %q: %w", listID, core_errors.ErrNotFound)
	}

	return nil
}
