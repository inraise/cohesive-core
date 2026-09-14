package budget_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *BudgetRepository) GetTotals(
	ctx context.Context,
	householdID uuid.UUID,
) (totalDeposits int64, totalExpenses int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE type = 'deposit'), 0) AS total_deposits,
			COALESCE(SUM(amount) FILTER (WHERE type = 'expense'), 0) AS total_expenses
		FROM household_transactions
		WHERE household_id = $1;
	`

	row := r.pool.QueryRow(ctx, query, householdID)

	if err := row.Scan(&totalDeposits, &totalExpenses); err != nil {
		return 0, 0, fmt.Errorf("scan error: %w", err)
	}

	return totalDeposits, totalExpenses, nil
}

func (r *BudgetRepository) GetMemberContributions(
	ctx context.Context,
	householdID uuid.UUID,
) ([]core_domain.MemberContribution, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT created_by, SUM(amount) AS total
		FROM household_transactions
		WHERE household_id = $1 AND type = 'deposit'
		GROUP BY created_by
		ORDER BY total DESC;
	`

	rows, err := r.pool.Query(ctx, query, householdID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	result := make([]core_domain.MemberContribution, 0)

	for rows.Next() {
		var contribution core_domain.MemberContribution
		if err := rows.Scan(&contribution.UserID, &contribution.Amount); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		result = append(result, contribution)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
