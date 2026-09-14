package budget_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	Type        string  `json:"type" validate:"required,oneof=deposit expense"`
	Amount      int64   `json:"amount" validate:"required,gt=0"`
	Description *string `json:"description" validate:"omitempty,max=300"`
}

func (s *BudgetService) CreateTransaction(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	request CreateTransactionRequest,
) (core_domain.HouseholdTransaction, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.HouseholdTransaction{}, err
	}

	txDomain := core_domain.NewHouseholdTransactionUninitialized(
		householdID,
		core_domain.TransactionType(request.Type),
		request.Amount,
		request.Description,
		callerID,
	)

	if err := txDomain.Validate(); err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("validate transaction domain: %w", err)
	}

	created, err := s.budgetRepository.CreateTransaction(ctx, txDomain)
	if err != nil {
		return core_domain.HouseholdTransaction{}, fmt.Errorf("create transaction: %w", err)
	}

	return created, nil
}
