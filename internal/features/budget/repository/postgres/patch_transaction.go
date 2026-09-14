package budget_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *BudgetRepository) PatchTransaction(
	ctx context.Context,
	tx core_domain.HouseholdTransaction,
) (core_domain.HouseholdTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE household_transactions
		SET amount = $1, description = $2, updated_at = $3, version = version + 1
		WHERE id = $4 AND household_id = $5 AND version = $6
		RETURNING id, household_id, version, type, amount, description, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		tx.Amount,
		tx.Description,
		time.Now(),
		tx.ID,
		tx.HouseholdID,
		tx.Version,
	)

	updated, err := scanTransaction(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.HouseholdTransaction{}, fmt.Errorf(
				"transaction with id %q concurrently accessed: %w", tx.ID, core_errors.ErrConflict,
			)
		}

		return core_domain.HouseholdTransaction{}, err
	}

	return updated, nil
}
