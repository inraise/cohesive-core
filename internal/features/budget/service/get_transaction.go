package budget_service

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

func (s *BudgetService) GetTransaction(
	ctx context.Context,
	householdID uuid.UUID,
	txID uuid.UUID,
	callerID uuid.UUID,
) (core_domain.HouseholdTransaction, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.HouseholdTransaction{}, err
	}

	tx, err := s.budgetRepository.GetTransactionByID(ctx, householdID, txID)
	if err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("get transaction from repository: %w", err)
	}

	return tx, nil
}
