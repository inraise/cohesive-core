package budget_repository_postgres

import (
	"context"

	core_domain "cohesive-core/internal/core/domain"
)

func (r *BudgetRepository) CreateTransaction(
	ctx context.Context,
	tx core_domain.HouseholdTransaction,
) (core_domain.HouseholdTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO household_transactions (household_id, type, amount, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, household_id, version, type, amount, description, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		tx.HouseholdID,
		string(tx.Type),
		tx.Amount,
		tx.Description,
		tx.CreatedBy,
		tx.CreatedAt,
		tx.UpdatedAt,
	)

	return scanTransaction(row)
}
