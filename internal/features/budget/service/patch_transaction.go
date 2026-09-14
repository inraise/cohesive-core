package budget_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *BudgetService) PatchTransaction(
	ctx context.Context,
	householdID uuid.UUID,
	txID uuid.UUID,
	callerID uuid.UUID,
	patch core_domain.HouseholdTransactionPatch,
) (core_domain.HouseholdTransaction, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.HouseholdTransaction{}, err
	}

	tx, err := s.budgetRepository.GetTransactionByID(ctx, householdID, txID)
	if err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("get transaction from repository: %w", err)
	}

	if err := s.requireCanModify(ctx, householdID, callerID, tx); err != nil {
		return core_domain.HouseholdTransaction{}, err
	}

	if err := tx.ApplyPatch(patch); err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("apply transaction patch: %w", err)
	}

	updated, err := s.budgetRepository.PatchTransaction(ctx, tx)
	if err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("patch transaction: %w", err)
	}

	return updated, nil
}
