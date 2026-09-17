package budget_repository_postgres

import (
	"context"
	"fmt"

	core_errors "cohesive-core/internal/core/errors"

	"github.com/google/uuid"
)

func (r *BudgetRepository) DeleteTransaction(
	ctx context.Context,
	householdID uuid.UUID,
	txID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM household_transactions WHERE id = $1 AND household_id = $2;`

	tag, err := r.pool.Exec(ctx, query, txID, householdID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("transaction with id %q: %w", txID, core_errors.ErrNotFound)
	}

	return nil
}
