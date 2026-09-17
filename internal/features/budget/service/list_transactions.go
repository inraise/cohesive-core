package budget_service

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

func (s *BudgetService) ListTransactions(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) ([]core_domain.HouseholdTransaction, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return nil, err
	}

	transactions, err := s.budgetRepository.ListTransactions(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("list transactions from repository: %w", err)
	}

	return transactions, nil
}
