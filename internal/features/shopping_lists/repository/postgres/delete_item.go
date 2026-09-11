package shoppinglists_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *ShoppingListsRepository) DeleteItem(
	ctx context.Context,
	shoppingListID uuid.UUID,
	itemID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM shopping_list_items WHERE id = $1 AND shopping_list_id = $2;`

	tag, err := r.pool.Exec(ctx, query, itemID, shoppingListID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item with id %q: %w", itemID, core_errors.ErrNotFound)
	}

	return nil
}
