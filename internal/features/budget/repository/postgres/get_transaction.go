package budget_repository_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"

	"github.com/google/uuid"
)

func (r *BudgetRepository) GetTransactionByID(
	ctx context.Context,
	householdID uuid.UUID,
	txID uuid.UUID,
) (core_domain.HouseholdTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, household_id, version, type, amount, description, created_by, created_at, updated_at
		FROM household_transactions
		WHERE id = $1 AND household_id = $2;
	`

	row := r.pool.QueryRow(ctx, query, txID, householdID)

	tx, err := scanTransaction(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.HouseholdTransaction{}, fmt.Errorf(
				"transaction with id %q in household %q: %w", txID, householdID, core_errors.ErrNotFound,
			)
		}

		return core_domain.HouseholdTransaction{}, err
	}

	return tx, nil
}
