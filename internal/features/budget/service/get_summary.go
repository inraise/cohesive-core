package budget_service

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"

	"github.com/google/uuid"
)

type BudgetSummary struct {
	Balance       int64
	TotalDeposits int64
	TotalExpenses int64
	Contributions []core_domain.MemberContribution
}

func (s *BudgetService) GetSummary(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) (BudgetSummary, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return BudgetSummary{}, err
	}

	totalDeposits, totalExpenses, err := s.budgetRepository.GetTotals(ctx, householdID)
	if err != nil {
		return BudgetSummary{}, fmt.Errorf("get totals from repository: %w", err)
	}

	contributions, err := s.budgetRepository.GetMemberContributions(ctx, householdID)
	if err != nil {
		return BudgetSummary{}, fmt.Errorf("get member contributions from repository: %w", err)
	}

	return BudgetSummary{
		Balance:       totalDeposits - totalExpenses,
		TotalDeposits: totalDeposits,
		TotalExpenses: totalExpenses,
		Contributions: contributions,
	}, nil
}
