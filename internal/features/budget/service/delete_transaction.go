package budget_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *BudgetService) DeleteTransaction(
	ctx context.Context,
	householdID uuid.UUID,
	txID uuid.UUID,
	callerID uuid.UUID,
) error {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return err
	}

	tx, err := s.budgetRepository.GetTransactionByID(ctx, householdID, txID)
	if err != nil {
		return fmt.Errorf("get transaction from repository: %w", err)
	}

	if err := s.requireCanModify(ctx, householdID, callerID, tx); err != nil {
		return err
	}

	if err := s.budgetRepository.DeleteTransaction(ctx, householdID, txID); err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	return nil
}
